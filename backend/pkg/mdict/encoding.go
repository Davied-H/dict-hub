package mdict

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// DecodeUTF16LE decodes UTF-16 Little Endian bytes to a UTF-8 string.
func DecodeUTF16LE(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// DecodeUTF16LEToString is a fast UTF-16LE decoder that converts bytes to string.
// It handles the conversion without relying on external transform.
func DecodeUTF16LEToString(data []byte, offset, length int) string {
	if length == 0 || offset+length > len(data) {
		return ""
	}
	
	src := data[offset : offset+length]
	
	// Handle odd length by truncating
	if len(src)%2 != 0 {
		src = src[:len(src)-1]
	}
	
	u16s := make([]uint16, len(src)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = uint16(src[i*2]) | uint16(src[i*2+1])<<8
	}
	
	// Remove null terminator if present
	for len(u16s) > 0 && u16s[len(u16s)-1] == 0 {
		u16s = u16s[:len(u16s)-1]
	}
	
	return string(utf16.Decode(u16s))
}

// DecodeUTF16BE decodes UTF-16 Big Endian bytes to a UTF-8 string.
func DecodeUTF16BE(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	
	decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// DecodeGBK decodes GBK encoded bytes to a UTF-8 string.
func DecodeGBK(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	
	decoder := simplifiedchinese.GBK.NewDecoder()
	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// DecodeGB18030 decodes GB18030 encoded bytes to a UTF-8 string.
func DecodeGB18030(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	
	decoder := simplifiedchinese.GB18030.NewDecoder()
	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// DecodeBig5 decodes Big5 encoded bytes to a UTF-8 string.
func DecodeBig5(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	
	decoder := traditionalchinese.Big5.NewDecoder()
	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// DecodeByEncoding decodes bytes based on the specified encoding.
func DecodeByEncoding(data []byte, enc Encoding) (string, error) {
	switch enc {
	case EncodingUTF16:
		return DecodeUTF16LE(data)
	case EncodingGBK, EncodingGB2312:
		return DecodeGBK(data)
	case EncodingGB18030:
		return DecodeGB18030(data)
	case EncodingBig5:
		return DecodeBig5(data)
	default:
		return string(data), nil
	}
}

// ReadBigEndianU64 reads a big-endian uint64 from bytes.
func ReadBigEndianU64(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(data)
}

// ReadBigEndianU32 reads a big-endian uint32 from bytes.
func ReadBigEndianU32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return binary.BigEndian.Uint32(data)
}

// ReadBigEndianU16 reads a big-endian uint16 from bytes.
func ReadBigEndianU16(data []byte) uint16 {
	if len(data) < 2 {
		return 0
	}
	return binary.BigEndian.Uint16(data)
}

// ReadLittleEndianU32 reads a little-endian uint32 from bytes.
func ReadLittleEndianU32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(data)
}

// ReadNumber reads a number from bytes based on the width (4 or 8 bytes).
func ReadNumber(data []byte, width int) int64 {
	if width == 8 {
		return int64(ReadBigEndianU64(data))
	}
	return int64(ReadBigEndianU32(data))
}

// NullTerminatedString extracts a null-terminated string from bytes.
func NullTerminatedString(data []byte, encoding Encoding, isMDD bool) string {
	// For UTF-16 or MDD files, look for double null bytes
	if encoding == EncodingUTF16 || isMDD {
		for i := 0; i < len(data)-1; i += 2 {
			if data[i] == 0 && data[i+1] == 0 {
				str, _ := DecodeUTF16LE(data[:i])
				return str
			}
		}
		str, _ := DecodeUTF16LE(data)
		return str
	}
	
	// For other encodings, look for single null byte
	idx := bytes.IndexByte(data, 0)
	if idx >= 0 {
		data = data[:idx]
	}
	
	str, _ := DecodeByEncoding(data, encoding)
	return str
}
