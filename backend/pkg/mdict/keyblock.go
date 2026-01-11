package mdict

import (
	"bytes"
	"fmt"
	"hash/adler32"
	"os"
)

// ReadKeyBlockMeta reads the key block metadata section.
// This is the FIXED implementation that correctly handles type 2 encryption.
//
// IMPORTANT: For type 2 encryption, only the Key Block INFO is encrypted,
// NOT the Key Block Metadata! The original library had this wrong.
func ReadKeyBlockMeta(file *os.File, header *Header) (*KeyBlockMeta, error) {
	meta := &KeyBlockMeta{}
	
	// Determine metadata size based on version
	var metaSize int64
	if header.Version >= 2.0 {
		metaSize = 40 // 5 x 8 bytes
	} else {
		metaSize = 16 // 4 x 4 bytes
	}
	
	// Read metadata from file
	data, err := ReadFileSection(file, header.HeaderEndPos, metaSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read key block metadata: %w", err)
	}
	
	// NOTE: Key Block Metadata is NOT encrypted for type 2 encryption.
	// Only the Key Block INFO section is encrypted.
	
	// Parse metadata based on version
	offset := 0
	width := header.NumberWidth
	
	// 1. Number of key blocks
	meta.KeyBlockNum = ReadNumber(data[offset:], width)
	offset += width
	
	// 2. Number of entries
	meta.EntriesNum = ReadNumber(data[offset:], width)
	offset += width
	
	// 3. Decompressed size of key block info (v2.0+ only)
	if header.Version >= 2.0 {
		meta.KeyBlockInfoDecompSize = ReadNumber(data[offset:], width)
		offset += width
	}
	
	// 4. Compressed size of key block info
	meta.KeyBlockInfoCompSize = ReadNumber(data[offset:], width)
	offset += width
	
	// 5. Total size of key blocks
	meta.KeyBlocksTotalSize = ReadNumber(data[offset:], width)
	
	// Calculate key block info start position
	// For v2.0+: header + 40 bytes metadata + 4 bytes checksum
	// For v1.x: header + 16 bytes metadata
	if header.Version >= 2.0 {
		meta.KeyBlockInfoStartPos = header.HeaderEndPos + 40 + 4
	} else {
		meta.KeyBlockInfoStartPos = header.HeaderEndPos + 16
	}
	
	return meta, nil
}

// ReadKeyBlockInfo reads and parses the key block info section.
// This contains metadata about each individual key block.
func ReadKeyBlockInfo(file *os.File, header *Header, meta *KeyBlockMeta) ([]*KeyBlockInfo, error) {
	// Read the compressed key block info
	data, err := ReadFileSection(file, meta.KeyBlockInfoStartPos, meta.KeyBlockInfoCompSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read key block info: %w", err)
	}
	
	// CRITICAL FIX: Decrypt key block info if encryption type is 2
	// This uses a different key derivation than metadata!
	if header.EncryptType == EncryptKeyInfo {
		data = DecryptKeyBlockInfo(data)
	}
	
	// Decompress the key block info
	decompressedData, _, err := DecompressBlock(data, meta.KeyBlockInfoDecompSize)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress key block info: %w", err)
	}
	
	// Verify decompressed size
	if int64(len(decompressedData)) != meta.KeyBlockInfoDecompSize {
		return nil, fmt.Errorf("decompressed size mismatch: expected %d, got %d",
			meta.KeyBlockInfoDecompSize, len(decompressedData))
	}
	
	// Verify checksum
	expectedChecksum := GetBlockChecksum(data)
	actualChecksum := adler32.Checksum(decompressedData)
	if actualChecksum != expectedChecksum {
		return nil, fmt.Errorf("key block info checksum mismatch: expected %d, got %d",
			expectedChecksum, actualChecksum)
	}
	
	// Parse key block info entries
	return parseKeyBlockInfoEntries(decompressedData, header, meta)
}

// parseKeyBlockInfoEntries parses the decompressed key block info into individual entries.
func parseKeyBlockInfoEntries(data []byte, header *Header, meta *KeyBlockMeta) ([]*KeyBlockInfo, error) {
	infos := make([]*KeyBlockInfo, 0, meta.KeyBlockNum)
	
	offset := 0
	width := header.NumberWidth
	
	// Determine terminator size based on version
	textTermSize := 0
	if header.Version >= 2.0 {
		textTermSize = 1
	}
	
	// Determine if using UTF-16 for keys
	isUTF16 := header.Encoding == EncodingUTF16
	
	var compAccum, decompAccum int64
	
	for i := int64(0); i < meta.KeyBlockNum; i++ {
		info := &KeyBlockInfo{}
		
		// Read number of entries in this block (not used directly but must be read)
		_ = ReadNumber(data[offset:], width)
		offset += width
		
		// Read first key size
		var firstKeySize int
		if header.Version >= 2.0 {
			firstKeySize = int(ReadBigEndianU16(data[offset:]))
			offset += 2
		} else {
			firstKeySize = int(data[offset])
			offset += 1
		}
		
		// Calculate actual byte size for first key
		firstKeyByteSize := firstKeySize
		if isUTF16 {
			firstKeyByteSize = (firstKeySize + textTermSize) * 2
		} else {
			firstKeyByteSize = firstKeySize + textTermSize
		}
		
		// Read first key
		if isUTF16 {
			info.FirstKey = DecodeUTF16LEToString(data, offset, firstKeyByteSize-textTermSize*2)
		} else {
			info.FirstKey = string(data[offset : offset+firstKeySize])
		}
		offset += firstKeyByteSize
		
		// Read last key size
		var lastKeySize int
		if header.Version >= 2.0 {
			lastKeySize = int(ReadBigEndianU16(data[offset:]))
			offset += 2
		} else {
			lastKeySize = int(data[offset])
			offset += 1
		}
		
		// Calculate actual byte size for last key
		lastKeyByteSize := lastKeySize
		if isUTF16 {
			lastKeyByteSize = (lastKeySize + textTermSize) * 2
		} else {
			lastKeyByteSize = lastKeySize + textTermSize
		}
		
		// Read last key
		if isUTF16 {
			info.LastKey = DecodeUTF16LEToString(data, offset, lastKeyByteSize-textTermSize*2)
		} else {
			info.LastKey = string(data[offset : offset+lastKeySize])
		}
		offset += lastKeyByteSize
		
		// Read compressed size
		info.CompressedSize = ReadNumber(data[offset:], width)
		offset += width
		
		// Read decompressed size
		info.DecompressedSize = ReadNumber(data[offset:], width)
		offset += width
		
		// Set accumulated offsets
		info.CompressedOffset = compAccum
		info.DecompressedOffset = decompAccum
		
		compAccum += info.CompressedSize
		decompAccum += info.DecompressedSize
		
		infos = append(infos, info)
	}
	
	// Verify total compressed size matches metadata
	if compAccum != meta.KeyBlocksTotalSize {
		return nil, fmt.Errorf("total key block size mismatch: expected %d, got %d",
			meta.KeyBlocksTotalSize, compAccum)
	}
	
	return infos, nil
}

// ReadKeyEntries reads all keyword entries from the key blocks.
func ReadKeyEntries(file *os.File, header *Header, meta *KeyBlockMeta, infos []*KeyBlockInfo, keyBlockDataStartPos int64) ([]*KeyEntry, error) {
	// Read all key block data
	data, err := ReadFileSection(file, keyBlockDataStartPos, meta.KeyBlocksTotalSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read key block data: %w", err)
	}
	
	entries := make([]*KeyEntry, 0, meta.EntriesNum)
	
	for _, info := range infos {
		// Get the compressed key block
		blockData := data[info.CompressedOffset : info.CompressedOffset+info.CompressedSize]
		
		// Decompress the key block
		decompressedBlock, _, err := DecompressBlock(blockData, info.DecompressedSize)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress key block: %w", err)
		}
		
		// Verify checksum
		expectedChecksum := GetBlockChecksum(blockData)
		actualChecksum := adler32.Checksum(decompressedBlock)
		if actualChecksum != expectedChecksum {
			return nil, fmt.Errorf("key block checksum mismatch: expected %d, got %d",
				expectedChecksum, actualChecksum)
		}
		
		// Parse entries from this block
		blockEntries := parseKeyBlockEntries(decompressedBlock, header)
		entries = append(entries, blockEntries...)
	}
	
	// Set record end offsets
	for i := 0; i < len(entries)-1; i++ {
		entries[i].RecordEndOffset = entries[i+1].RecordStartOffset
	}
	
	// Verify entry count
	if int64(len(entries)) != meta.EntriesNum {
		return nil, fmt.Errorf("entry count mismatch: expected %d, got %d",
			meta.EntriesNum, len(entries))
	}
	
	return entries, nil
}

// parseKeyBlockEntries parses entries from a decompressed key block.
func parseKeyBlockEntries(data []byte, header *Header) []*KeyEntry {
	entries := make([]*KeyEntry, 0)
	
	width := header.NumberWidth
	isUTF16 := header.Encoding == EncodingUTF16
	
	// Determine null terminator width
	termWidth := 1
	if isUTF16 {
		termWidth = 2
	}
	
	offset := 0
	for offset < len(data) {
		// Read record start offset
		recordOffset := ReadNumber(data[offset:], width)
		offset += width
		
		// Find null terminator
		keyEnd := offset
		if isUTF16 {
			// Look for double null byte
			for keyEnd < len(data)-1 {
				if data[keyEnd] == 0 && data[keyEnd+1] == 0 {
					break
				}
				keyEnd += 2
			}
		} else {
			// Look for single null byte
			keyEnd = bytes.IndexByte(data[offset:], 0)
			if keyEnd < 0 {
				keyEnd = len(data)
			} else {
				keyEnd += offset
			}
		}
		
		// Extract keyword
		keyBytes := data[offset:keyEnd]
		var keyword string
		if isUTF16 {
			keyword, _ = DecodeUTF16LE(keyBytes)
		} else {
			keyword, _ = DecodeByEncoding(keyBytes, header.Encoding)
		}
		
		entries = append(entries, &KeyEntry{
			Keyword:           keyword,
			RecordStartOffset: recordOffset,
		})
		
		// Move past the null terminator
		offset = keyEnd + termWidth
	}
	
	return entries
}
