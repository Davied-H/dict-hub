// Package mdict provides MDX/MDD dictionary file parsing functionality.
// This implementation is based on the Python readmdict library.
package mdict

// DictType represents the type of dictionary file.
type DictType int

const (
	// DictTypeMDX indicates an MDX file (dictionary definitions).
	DictTypeMDX DictType = iota
	// DictTypeMDD indicates an MDD file (resources like images, audio).
	DictTypeMDD
)

// EncryptType represents the encryption type used in the dictionary.
type EncryptType int

const (
	// EncryptNone indicates no encryption.
	EncryptNone EncryptType = 0
	// EncryptRecord indicates record block encryption.
	EncryptRecord EncryptType = 1
	// EncryptKeyInfo indicates key info block encryption.
	EncryptKeyInfo EncryptType = 2
)

// Encoding represents the text encoding used in the dictionary.
type Encoding int

const (
	EncodingUTF8 Encoding = iota
	EncodingUTF16
	EncodingGBK
	EncodingGB2312
	EncodingGB18030
	EncodingBig5
)

// CompressionType represents the compression algorithm used.
type CompressionType byte

const (
	CompressionNone CompressionType = 0
	CompressionLZO  CompressionType = 1
	CompressionZlib CompressionType = 2
)

// Header contains the parsed header information from the dictionary file.
type Header struct {
	// Raw header data
	HeaderSize    uint32
	HeaderBytes   []byte
	HeaderXML     string
	Adler32       uint32
	HeaderEndPos  int64

	// Parsed metadata from XML
	Title                    string
	Description              string
	CreationDate             string
	GeneratedByEngineVersion string
	DataSourceFormat         string
	StyleSheet               string
	
	// Parsed settings
	Version     float64
	Encoding    Encoding
	EncryptType EncryptType
	
	// Computed values based on version
	NumberWidth int // 4 for v1.x, 8 for v2.x
}

// KeyBlockMeta contains metadata about key blocks.
type KeyBlockMeta struct {
	KeyBlockNum            int64 // Number of key blocks
	EntriesNum             int64 // Total number of entries
	KeyBlockInfoDecompSize int64 // Decompressed size of key block info (v2.0 only)
	KeyBlockInfoCompSize   int64 // Compressed size of key block info
	KeyBlocksTotalSize     int64 // Total size of all key blocks
	
	// File positions
	KeyBlockInfoStartPos int64
}

// KeyBlockInfo contains information about a single key block.
type KeyBlockInfo struct {
	FirstKey         string
	LastKey          string
	CompressedSize   int64
	DecompressedSize int64
	
	// Accumulated offsets for quick access
	CompressedOffset   int64
	DecompressedOffset int64
}

// KeyEntry represents a single keyword entry.
type KeyEntry struct {
	Keyword           string
	RecordStartOffset int64
	RecordEndOffset   int64
}

// RecordBlockMeta contains metadata about record blocks.
type RecordBlockMeta struct {
	RecordBlockNum        int64 // Number of record blocks
	EntriesNum            int64 // Total number of entries (should match KeyBlockMeta.EntriesNum)
	RecordBlockInfoSize   int64 // Size of record block info section
	RecordBlocksTotalSize int64 // Total size of all record blocks
	
	// File positions
	RecordBlockMetaStartPos int64
	RecordBlockMetaEndPos   int64
}

// RecordBlockInfo contains information about a single record block.
type RecordBlockInfo struct {
	CompressedSize   int64
	DecompressedSize int64
	
	// Accumulated offsets
	CompressedOffset   int64
	DecompressedOffset int64
}

// Mdict represents a parsed MDX/MDD dictionary.
type Mdict struct {
	FilePath string
	DictType DictType
	
	Header          *Header
	KeyBlockMeta    *KeyBlockMeta
	KeyBlockInfos   []*KeyBlockInfo
	KeyEntries      []*KeyEntry
	RecordBlockMeta *RecordBlockMeta
	RecordBlockInfos []*RecordBlockInfo
	
	// File positions
	KeyBlockDataStartPos    int64
	RecordBlockDataStartPos int64
}

// DictInfo contains basic dictionary information for API responses.
type DictInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
	HasMDD      bool   `json:"has_mdd"`
	WordCount   int64  `json:"word_count"`
}

// SearchResult represents a search result from the dictionary.
type SearchResult struct {
	DictID     uint   `json:"dict_id"`
	DictName   string `json:"dict_name"`
	DictTitle  string `json:"dict_title"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

// SuggestResult represents a suggestion result.
type SuggestResult struct {
	Word      string `json:"word"`
	DictID    uint   `json:"dict_id"`
	DictTitle string `json:"dict_title"`
}
