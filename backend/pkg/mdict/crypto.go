package mdict

import (
	"github.com/c0mm4nd/go-ripemd"
)

// ripemd128 computes RIPEMD-128 hash of the input data.
func ripemd128(data []byte) []byte {
	h := ripemd.New128()
	h.Write(data)
	return h.Sum(nil)
}

// fastDecrypt performs the MDict fast decryption algorithm.
// This is used for decrypting key block info and record blocks.
//
// The algorithm:
// 1. For each byte, swap nibbles (high 4 bits <-> low 4 bits)
// 2. XOR with previous byte
// 3. XOR with position (i mod 256)
// 4. XOR with key byte (key[i mod keyLen])
func fastDecrypt(data []byte, key []byte) {
	previous := byte(0x36)
	keyLen := len(key)
	
	for i := 0; i < len(data); i++ {
		// Swap nibbles: rotate byte by 4 bits
		t := ((data[i] >> 4) | (data[i] << 4)) & 0xff
		// XOR chain
		t ^= previous
		t ^= byte(i & 0xff)
		t ^= key[i%keyLen]
		// Save original for next iteration
		previous = data[i]
		data[i] = t
	}
}

// DecryptKeyBlockInfo decrypts the key block info section when encryption type is 2.
// This is the FIXED implementation based on Python readmdict.
//
// The key is derived from:
// 1. Take first 4 bytes of the compressed data (which contains adler32 checksum)
// 2. Append magic bytes 0x95, 0x36, 0x00, 0x00
// 3. Hash with RIPEMD-128 to get the decryption key
func DecryptKeyBlockInfo(data []byte) []byte {
	if len(data) < 8 {
		return data
	}
	
	// Build the key buffer:
	// [4 bytes from data[4:8]] + [0x95, 0x36, 0x00, 0x00]
	keyBuffer := make([]byte, 8)
	copy(keyBuffer[0:4], data[4:8]) // adler32 checksum bytes
	keyBuffer[4] = 0x95
	keyBuffer[5] = 0x36
	keyBuffer[6] = 0x00
	keyBuffer[7] = 0x00
	
	// Generate decryption key using RIPEMD-128
	key := ripemd128(keyBuffer)
	
	// Decrypt in place starting from offset 8
	// (skip the 4-byte compression type and 4-byte checksum)
	fastDecrypt(data[8:], key)
	
	return data
}

// DecryptRecordBlock decrypts a record block when encryption type is 1.
// Similar to key block info decryption but applied to record data.
func DecryptRecordBlock(data []byte, compressedSize int64) []byte {
	if len(data) < 8 {
		return data
	}
	
	// Build the key buffer from the checksum
	keyBuffer := make([]byte, 8)
	copy(keyBuffer[0:4], data[4:8])
	keyBuffer[4] = 0x95
	keyBuffer[5] = 0x36
	keyBuffer[6] = 0x00
	keyBuffer[7] = 0x00
	
	key := ripemd128(keyBuffer)
	
	// Decrypt data after the 8-byte header
	fastDecrypt(data[8:], key)
	
	return data
}

// DecryptKeyBlockMeta decrypts the key block metadata when encryption type is 2.
// For version 2.0+, the metadata is 40 bytes (5 x 8-byte values).
// For version 1.x, the metadata is 16 bytes (4 x 4-byte values).
//
// IMPORTANT: The key for metadata decryption uses the HEADER's adler32 checksum,
// not the data's own checksum.
func DecryptKeyBlockMeta(data []byte, headerAdler32 uint32) []byte {
	if len(data) < 8 {
		return data
	}
	
	// Build key from header's adler32 checksum
	keyBuffer := make([]byte, 8)
	// Convert uint32 to big-endian bytes
	keyBuffer[0] = byte(headerAdler32 >> 24)
	keyBuffer[1] = byte(headerAdler32 >> 16)
	keyBuffer[2] = byte(headerAdler32 >> 8)
	keyBuffer[3] = byte(headerAdler32)
	keyBuffer[4] = 0x95
	keyBuffer[5] = 0x36
	keyBuffer[6] = 0x00
	keyBuffer[7] = 0x00
	
	key := ripemd128(keyBuffer)
	
	// Decrypt the entire metadata buffer
	fastDecrypt(data, key)
	
	return data
}
