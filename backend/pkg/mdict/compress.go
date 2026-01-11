package mdict

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/rasky/go-lzo"
)

// DecompressZlib decompresses zlib-compressed data.
func DecompressZlib(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data for zlib decompression")
	}
	
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()
	
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress zlib data: %w", err)
	}
	
	return result, nil
}

// DecompressLZO decompresses LZO-compressed data.
// expectedSize is the expected decompressed size (required for LZO).
func DecompressLZO(data []byte, expectedSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data for LZO decompression")
	}
	
	reader := bytes.NewReader(data)
	result, err := lzo.Decompress1X(reader, 0, expectedSize)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress LZO data: %w", err)
	}
	
	return result, nil
}

// DecompressBlock decompresses a data block based on the compression type.
// The first 4 bytes indicate the compression type:
// - 0x00: No compression
// - 0x01: LZO compression
// - 0x02: zlib compression
// The next 4 bytes are the adler32 checksum of the decompressed data.
// The remaining bytes are the compressed data.
func DecompressBlock(data []byte, expectedDecompSize int64) ([]byte, CompressionType, error) {
	if len(data) < 8 {
		return nil, CompressionNone, fmt.Errorf("data block too small: %d bytes", len(data))
	}
	
	compressionType := CompressionType(data[0])
	// data[1:4] should be zeros in standard format
	// data[4:8] is the adler32 checksum
	compressedData := data[8:]
	
	var result []byte
	var err error
	
	switch compressionType {
	case CompressionNone:
		// No compression, data is raw
		result = compressedData
		
	case CompressionLZO:
		result, err = DecompressLZO(compressedData, int(expectedDecompSize))
		if err != nil {
			return nil, compressionType, fmt.Errorf("LZO decompression failed: %w", err)
		}
		
	case CompressionZlib:
		result, err = DecompressZlib(compressedData)
		if err != nil {
			return nil, compressionType, fmt.Errorf("zlib decompression failed: %w", err)
		}
		
	default:
		return nil, compressionType, fmt.Errorf("unknown compression type: %d", compressionType)
	}
	
	return result, compressionType, nil
}

// GetCompressionType returns the compression type from a data block's header.
func GetCompressionType(data []byte) CompressionType {
	if len(data) < 1 {
		return CompressionNone
	}
	return CompressionType(data[0])
}

// GetBlockChecksum returns the adler32 checksum from a data block's header.
func GetBlockChecksum(data []byte) uint32 {
	if len(data) < 8 {
		return 0
	}
	return ReadBigEndianU32(data[4:8])
}
