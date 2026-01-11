package mdict

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// New creates a new Mdict instance from a file path.
// It reads and parses the header and key block metadata.
func New(filePath string) (*Mdict, error) {
	// Determine dictionary type from extension
	dictType := DictTypeMDX
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".mdd" {
		dictType = DictTypeMDD
	}
	
	mdict := &Mdict{
		FilePath: filePath,
		DictType: dictType,
	}
	
	// Open and parse the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Read header
	header, err := ReadHeader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	mdict.Header = header
	
	// Correct encoding for MDD files
	if dictType == DictTypeMDD {
		mdict.Header.Encoding = EncodingUTF16
	}
	
	// Read key block metadata
	keyBlockMeta, err := ReadKeyBlockMeta(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to read key block metadata: %w", err)
	}
	mdict.KeyBlockMeta = keyBlockMeta
	
	return mdict, nil
}

// BuildIndex builds the full dictionary index by reading all key blocks and record blocks.
// This should be called after New() before using Lookup().
func (m *Mdict) BuildIndex() error {
	file, err := os.Open(m.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file for indexing: %w", err)
	}
	defer file.Close()
	
	// Read key block info
	keyBlockInfos, err := ReadKeyBlockInfo(file, m.Header, m.KeyBlockMeta)
	if err != nil {
		return fmt.Errorf("failed to read key block info: %w", err)
	}
	m.KeyBlockInfos = keyBlockInfos
	
	// Calculate key block data start position
	m.KeyBlockDataStartPos = m.KeyBlockMeta.KeyBlockInfoStartPos + m.KeyBlockMeta.KeyBlockInfoCompSize
	
	// Read key entries
	keyEntries, err := ReadKeyEntries(file, m.Header, m.KeyBlockMeta, keyBlockInfos, m.KeyBlockDataStartPos)
	if err != nil {
		return fmt.Errorf("failed to read key entries: %w", err)
	}
	m.KeyEntries = keyEntries
	
	// Calculate record block metadata start position
	recordBlockMetaStartPos := m.KeyBlockDataStartPos + m.KeyBlockMeta.KeyBlocksTotalSize
	
	// Read record block metadata
	recordBlockMeta, err := ReadRecordBlockMeta(file, m.Header, recordBlockMetaStartPos)
	if err != nil {
		return fmt.Errorf("failed to read record block metadata: %w", err)
	}
	m.RecordBlockMeta = recordBlockMeta
	
	// Verify entry counts match
	if recordBlockMeta.EntriesNum != m.KeyBlockMeta.EntriesNum {
		return fmt.Errorf("entry count mismatch: key blocks have %d, record blocks have %d",
			m.KeyBlockMeta.EntriesNum, recordBlockMeta.EntriesNum)
	}
	
	// Read record block info
	recordBlockInfos, err := ReadRecordBlockInfo(file, m.Header, recordBlockMeta)
	if err != nil {
		return fmt.Errorf("failed to read record block info: %w", err)
	}
	m.RecordBlockInfos = recordBlockInfos
	
	// Calculate record block data start position
	m.RecordBlockDataStartPos = recordBlockMeta.RecordBlockMetaEndPos + recordBlockMeta.RecordBlockInfoSize
	
	return nil
}

// Lookup looks up a word in the dictionary and returns its definition.
// Uses binary search for efficient lookup.
func (m *Mdict) Lookup(word string) ([]byte, error) {
	if m.KeyEntries == nil || len(m.KeyEntries) == 0 {
		return nil, fmt.Errorf("dictionary index not built, call BuildIndex() first")
	}
	
	word = strings.TrimSpace(word)
	
	// Binary search for the word
	idx := sort.Search(len(m.KeyEntries), func(i int) bool {
		return strings.ToLower(m.KeyEntries[i].Keyword) >= strings.ToLower(word)
	})
	
	// Check if found
	if idx >= len(m.KeyEntries) || strings.ToLower(m.KeyEntries[idx].Keyword) != strings.ToLower(word) {
		return nil, fmt.Errorf("word not found: %s", word)
	}
	
	return m.LookupByEntry(m.KeyEntries[idx])
}

// LookupByEntry looks up a definition by its key entry.
func (m *Mdict) LookupByEntry(entry *KeyEntry) ([]byte, error) {
	file, err := os.Open(m.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	data, err := LookupRecord(file, m.Header, entry, m.RecordBlockInfos, m.RecordBlockDataStartPos)
	if err != nil {
		return nil, err
	}
	
	// For MDD files, return raw data
	if m.DictType == DictTypeMDD {
		return data, nil
	}
	
	// For MDX files with UTF-16 encoding, decode to UTF-8
	if m.Header.Encoding == EncodingUTF16 {
		str, err := DecodeUTF16LE(data)
		if err != nil {
			return data, nil // Return raw data on decode error
		}
		return []byte(str), nil
	}
	
	// For other encodings, decode appropriately
	str, err := DecodeByEncoding(data, m.Header.Encoding)
	if err != nil {
		return data, nil
	}
	return []byte(str), nil
}

// Name returns the dictionary name (filename without extension).
func (m *Mdict) Name() string {
	name := filepath.Base(m.FilePath)
	name = strings.TrimSuffix(name, ".mdx")
	name = strings.TrimSuffix(name, ".mdd")
	name = strings.TrimSuffix(name, ".MDX")
	name = strings.TrimSuffix(name, ".MDD")
	return name
}

// Title returns the dictionary title from the header.
func (m *Mdict) Title() string {
	if m.Header != nil && m.Header.Title != "" {
		return m.Header.Title
	}
	return m.Name()
}

// Description returns the dictionary description from the header.
func (m *Mdict) Description() string {
	if m.Header != nil {
		return m.Header.Description
	}
	return ""
}

// WordCount returns the total number of entries in the dictionary.
func (m *Mdict) WordCount() int64 {
	if m.KeyBlockMeta != nil {
		return m.KeyBlockMeta.EntriesNum
	}
	return 0
}

// IsMDD returns true if this is an MDD file (resources).
func (m *Mdict) IsMDD() bool {
	return m.DictType == DictTypeMDD
}

// GetKeyEntries returns all keyword entries.
func (m *Mdict) GetKeyEntries() []*KeyEntry {
	return m.KeyEntries
}

// Suggest returns word suggestions based on a prefix.
func (m *Mdict) Suggest(prefix string, limit int) []string {
	if m.KeyEntries == nil || len(m.KeyEntries) == 0 {
		return nil
	}
	
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return nil
	}
	
	// Find the first entry with matching prefix
	idx := sort.Search(len(m.KeyEntries), func(i int) bool {
		return strings.ToLower(m.KeyEntries[i].Keyword) >= prefix
	})
	
	// Collect matching entries
	results := make([]string, 0, limit)
	seen := make(map[string]bool)
	
	for i := idx; i < len(m.KeyEntries) && len(results) < limit; i++ {
		keyword := m.KeyEntries[i].Keyword
		if !strings.HasPrefix(strings.ToLower(keyword), prefix) {
			break
		}
		if !seen[keyword] {
			seen[keyword] = true
			results = append(results, keyword)
		}
	}
	
	return results
}

// Close releases any resources associated with the dictionary.
// Currently, this is a no-op as files are opened and closed for each operation.
func (m *Mdict) Close() error {
	return nil
}
