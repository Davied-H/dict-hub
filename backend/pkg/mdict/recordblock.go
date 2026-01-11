package mdict

import (
	"fmt"
	"hash/adler32"
	"os"
)

// ReadRecordBlockMeta reads the record block metadata section.
func ReadRecordBlockMeta(file *os.File, header *Header, recordBlockMetaStartPos int64) (*RecordBlockMeta, error) {
	meta := &RecordBlockMeta{
		RecordBlockMetaStartPos: recordBlockMetaStartPos,
	}
	
	// Determine metadata size based on version
	var metaSize int64
	if header.Version >= 2.0 {
		metaSize = 32 // 4 x 8 bytes
	} else {
		metaSize = 16 // 4 x 4 bytes
	}
	
	// Read metadata from file
	data, err := ReadFileSection(file, recordBlockMetaStartPos, metaSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read record block metadata: %w", err)
	}
	
	// Parse metadata
	offset := 0
	width := header.NumberWidth
	
	// 1. Number of record blocks
	meta.RecordBlockNum = ReadNumber(data[offset:], width)
	offset += width
	
	// 2. Number of entries
	meta.EntriesNum = ReadNumber(data[offset:], width)
	offset += width
	
	// 3. Size of record block info section
	meta.RecordBlockInfoSize = ReadNumber(data[offset:], width)
	offset += width
	
	// 4. Total size of record blocks
	meta.RecordBlocksTotalSize = ReadNumber(data[offset:], width)
	
	// Calculate metadata end position
	meta.RecordBlockMetaEndPos = recordBlockMetaStartPos + metaSize
	
	return meta, nil
}

// ReadRecordBlockInfo reads and parses the record block info section.
func ReadRecordBlockInfo(file *os.File, header *Header, meta *RecordBlockMeta) ([]*RecordBlockInfo, error) {
	// Read record block info data
	data, err := ReadFileSection(file, meta.RecordBlockMetaEndPos, meta.RecordBlockInfoSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read record block info: %w", err)
	}
	
	infos := make([]*RecordBlockInfo, 0, meta.RecordBlockNum)
	
	offset := 0
	width := header.NumberWidth
	var compAccum, decompAccum int64
	
	for i := int64(0); i < meta.RecordBlockNum; i++ {
		info := &RecordBlockInfo{}
		
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
	
	// Verify sizes
	if compAccum != meta.RecordBlocksTotalSize {
		return nil, fmt.Errorf("record block total size mismatch: expected %d, got %d",
			meta.RecordBlocksTotalSize, compAccum)
	}
	
	if int64(offset) != meta.RecordBlockInfoSize {
		return nil, fmt.Errorf("record block info size mismatch: expected %d, got %d",
			meta.RecordBlockInfoSize, offset)
	}
	
	return infos, nil
}

// LookupRecord looks up a word definition from the record blocks.
func LookupRecord(file *os.File, header *Header, entry *KeyEntry, infos []*RecordBlockInfo, recordBlockDataStartPos int64) ([]byte, error) {
	// Find the record block containing this entry
	var targetInfo *RecordBlockInfo
	for _, info := range infos {
		if entry.RecordStartOffset >= info.DecompressedOffset &&
			entry.RecordStartOffset < info.DecompressedOffset+info.DecompressedSize {
			targetInfo = info
			break
		}
	}
	
	if targetInfo == nil {
		return nil, fmt.Errorf("could not find record block for offset %d", entry.RecordStartOffset)
	}
	
	// Read the compressed record block
	blockData, err := ReadFileSection(file,
		recordBlockDataStartPos+targetInfo.CompressedOffset,
		targetInfo.CompressedSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read record block: %w", err)
	}
	
	// Decrypt if needed (type 1 encryption)
	if header.EncryptType == EncryptRecord {
		blockData = DecryptRecordBlock(blockData, targetInfo.CompressedSize)
	}
	
	// Decompress the record block
	decompressedBlock, _, err := DecompressBlock(blockData, targetInfo.DecompressedSize)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress record block: %w", err)
	}
	
	// Verify checksum
	expectedChecksum := GetBlockChecksum(blockData)
	actualChecksum := adler32.Checksum(decompressedBlock)
	if actualChecksum != expectedChecksum {
		return nil, fmt.Errorf("record block checksum mismatch: expected %d, got %d",
			expectedChecksum, actualChecksum)
	}
	
	// Calculate offsets within the decompressed block
	startOffset := entry.RecordStartOffset - targetInfo.DecompressedOffset
	var endOffset int64
	if entry.RecordEndOffset > 0 {
		endOffset = entry.RecordEndOffset - targetInfo.DecompressedOffset
	} else {
		endOffset = int64(len(decompressedBlock))
	}
	
	// Validate offsets
	if startOffset < 0 || startOffset >= int64(len(decompressedBlock)) {
		return nil, fmt.Errorf("invalid start offset: %d (block size: %d)",
			startOffset, len(decompressedBlock))
	}
	if endOffset < startOffset || endOffset > int64(len(decompressedBlock)) {
		endOffset = int64(len(decompressedBlock))
	}
	
	// Extract the definition data
	definitionData := decompressedBlock[startOffset:endOffset]
	
	return definitionData, nil
}

// FindRecordBlockForOffset finds the record block info containing the given decompressed offset.
func FindRecordBlockForOffset(infos []*RecordBlockInfo, offset int64) *RecordBlockInfo {
	for _, info := range infos {
		if offset >= info.DecompressedOffset && offset < info.DecompressedOffset+info.DecompressedSize {
			return info
		}
	}
	return nil
}
