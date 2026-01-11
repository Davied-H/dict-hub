package mdict

import (
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ReadHeader reads and parses the dictionary header from the file.
func ReadHeader(file *os.File) (*Header, error) {
	header := &Header{}
	
	// Read header size (4 bytes, big-endian)
	var headerSize uint32
	if err := binary.Read(file, binary.BigEndian, &headerSize); err != nil {
		return nil, fmt.Errorf("failed to read header size: %w", err)
	}
	header.HeaderSize = headerSize
	
	// Read header bytes
	headerBytes := make([]byte, headerSize)
	if _, err := io.ReadFull(file, headerBytes); err != nil {
		return nil, fmt.Errorf("failed to read header bytes: %w", err)
	}
	header.HeaderBytes = headerBytes
	
	// Read adler32 checksum (4 bytes, big-endian)
	if err := binary.Read(file, binary.BigEndian, &header.Adler32); err != nil {
		return nil, fmt.Errorf("failed to read header checksum: %w", err)
	}
	
	// Calculate header end position
	header.HeaderEndPos = int64(4 + headerSize + 4)
	
	// Decode header bytes (UTF-16LE encoded XML)
	headerXML := DecodeUTF16LEToString(headerBytes, 0, len(headerBytes))
	
	// Replace Library_Data with Dictionary for compatibility
	headerXML = strings.Replace(headerXML, "Library_Data", "Dictionary", 1)
	header.HeaderXML = headerXML
	
	// Verify checksum (optional, some dictionaries have mismatched checksums)
	calculatedChecksum := adler32.Checksum([]byte(headerXML))
	if calculatedChecksum != header.Adler32 {
		// Log warning but don't fail - some dictionaries have this issue
		// fmt.Printf("Warning: header checksum mismatch (expected %d, got %d)\n", header.Adler32, calculatedChecksum)
	}
	
	// Parse XML header info
	if err := parseHeaderXML(header); err != nil {
		return nil, fmt.Errorf("failed to parse header XML: %w", err)
	}
	
	return header, nil
}

// parseHeaderXML extracts metadata from the header XML string.
func parseHeaderXML(header *Header) error {
	xml := header.HeaderXML
	
	// Extract version
	header.Version = extractFloatAttr(xml, "GeneratedByEngineVersion", 2.0)
	header.GeneratedByEngineVersion = extractAttr(xml, "GeneratedByEngineVersion")
	
	// Set number width based on version
	if header.Version >= 2.0 {
		header.NumberWidth = 8
	} else {
		header.NumberWidth = 4
	}
	
	// Extract title
	header.Title = extractAttr(xml, "Title")
	if header.Title == "" {
		// Try to get title from between tags
		header.Title = extractTagContent(xml, "Title")
	}
	
	// Extract description
	header.Description = extractAttr(xml, "Description")
	if header.Description == "" {
		header.Description = extractTagContent(xml, "Description")
	}
	
	// Extract creation date
	header.CreationDate = extractAttr(xml, "CreationDate")
	
	// Extract data source format
	header.DataSourceFormat = extractAttr(xml, "DataSourceFormat")
	
	// Extract stylesheet
	header.StyleSheet = extractAttr(xml, "StyleSheet")
	if header.StyleSheet == "" {
		header.StyleSheet = extractTagContent(xml, "StyleSheet")
	}
	
	// Extract encryption type
	encrypted := extractAttr(xml, "Encrypted")
	header.EncryptType = parseEncryptType(encrypted)
	
	// Extract encoding
	encodingStr := extractAttr(xml, "Encoding")
	header.Encoding = parseEncoding(encodingStr)
	
	return nil
}

// extractAttr extracts an attribute value from XML-like string.
// Handles both attribute="value" and attribute=value formats.
func extractAttr(xml, attrName string) string {
	// Try attribute="value" format
	pattern := regexp.MustCompile(attrName + `="([^"]*)"`)
	matches := pattern.FindStringSubmatch(xml)
	if len(matches) >= 2 {
		return matches[1]
	}
	
	// Try attribute='value' format
	pattern = regexp.MustCompile(attrName + `='([^']*)'`)
	matches = pattern.FindStringSubmatch(xml)
	if len(matches) >= 2 {
		return matches[1]
	}
	
	return ""
}

// extractTagContent extracts content between XML tags.
func extractTagContent(xml, tagName string) string {
	pattern := regexp.MustCompile(`<` + tagName + `[^>]*>([\s\S]*?)</` + tagName + `>`)
	matches := pattern.FindStringSubmatch(xml)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractFloatAttr extracts a float attribute value.
func extractFloatAttr(xml, attrName string, defaultVal float64) float64 {
	str := extractAttr(xml, attrName)
	if str == "" {
		return defaultVal
	}
	
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// parseEncryptType parses the Encrypted attribute value.
func parseEncryptType(encrypted string) EncryptType {
	encrypted = strings.TrimSpace(encrypted)
	
	switch encrypted {
	case "", "No", "no", "NO", "false", "False", "FALSE", "0":
		return EncryptNone
	case "Yes", "yes", "YES", "true", "True", "TRUE", "1":
		return EncryptRecord
	case "2":
		return EncryptKeyInfo
	default:
		// Try to parse as number
		if len(encrypted) > 0 {
			if encrypted[0] == '2' {
				return EncryptKeyInfo
			}
			if encrypted[0] == '1' {
				return EncryptRecord
			}
		}
		return EncryptNone
	}
}

// parseEncoding parses the Encoding attribute value.
func parseEncoding(encoding string) Encoding {
	encoding = strings.ToLower(strings.TrimSpace(encoding))
	
	switch encoding {
	case "gbk":
		return EncodingGBK
	case "gb2312":
		return EncodingGB2312
	case "gb18030":
		return EncodingGB18030
	case "big5":
		return EncodingBig5
	case "utf-16", "utf16":
		return EncodingUTF16
	default:
		return EncodingUTF8
	}
}

// ReadFileSection reads a section of the file from a specific position.
func ReadFileSection(file *os.File, offset int64, length int64) ([]byte, error) {
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to position %d: %w", offset, err)
	}
	
	data := make([]byte, length)
	if _, err := io.ReadFull(file, data); err != nil {
		return nil, fmt.Errorf("failed to read %d bytes from position %d: %w", length, offset, err)
	}
	
	return data, nil
}
