//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mdx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"os"
	"strconv"
	"strings"

	"github.com/rasky/go-lzo"
)

// readDictHeader reads and parses the dictionary's header information.
// It first calls readMDictFileHeader to read the raw header data from the file,
// then validates the checksum, and finally parses the XML-formatted header info to populate the meta struct.
func (mdict *MdictBase) readDictHeader() error {
	log.Infof("Reading dictionary header: %s", mdict.filePath)
	dictHeader, err := readMDictFileHeader(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to read MDict file header for '%s': %w", mdict.filePath, err)
	}
	mdict.header = dictHeader

	// Note: The checksum is calculated on the headerInfo string after UTF-8 conversion and string replacement.
	// This is to maintain compatibility with the behavior of some MDict generation tools.
	log.Debugf("Verifying header checksum for '%s'. Expected: %d", mdict.filePath, dictHeader.adler32Checksum)
	checksum := adler32.Checksum([]byte(dictHeader.headerInfo))
	if checksum != dictHeader.adler32Checksum {
		// In some dictionaries, this checksum may not match, but the dictionary can still be parsed.
		// Therefore, we only log the error without interrupting the parsing.
		log.Warningf("Header checksum mismatch for '%s': expected %d, calculated %d", mdict.filePath, dictHeader.adler32Checksum, checksum)
	}
	log.Debugf("Header checksum verification complete for: %s", mdict.filePath)

	log.Debugf("Parsing header XML info for: %s", mdict.filePath)
	headerInfo, err := parseXMLHeader(dictHeader.headerInfo)
	if err != nil {
		return fmt.Errorf("failed to parse XML header for '%s': %w", mdict.filePath, err)
	}
	log.Debugf("Header info parsed for '%s'. Title: '%s', EngineVersion: '%s', Encoding: '%s'", mdict.filePath, headerInfo.Title, headerInfo.GeneratedByEngineVersion, headerInfo.Encoding)

	meta := &mdictMeta{}

	// Process encryption flag
	encrypted := headerInfo.Encrypted
	switch {
	case encrypted == "" || encrypted == "No":
		meta.encryptType = EncryptNoEnc
	case encrypted == "Yes":
		meta.encryptType = EncryptRecordEnc
	default:
		if len(encrypted) > 0 && encrypted[0] == '2' {
			meta.encryptType = EncryptKeyInfoEnc
		} else if len(encrypted) > 0 && encrypted[0] == '1' {
			meta.encryptType = EncryptRecordEnc
		} else {
			meta.encryptType = EncryptNoEnc
		}
	}

	// Process version number
	versionStr := headerInfo.GeneratedByEngineVersion
	version, err := strconv.ParseFloat(versionStr, 32)
	if err != nil {
		log.Errorf("Invalid engine version '%s' in header for '%s': %v", versionStr, mdict.filePath, err)
		return fmt.Errorf("invalid engine version '%s' in header for '%s': %w", versionStr, mdict.filePath, err)
	}
	meta.version = float32(version)
	log.Debugf("Mdict version for '%s': %.1f", mdict.filePath, meta.version)

	// Process number format and width based on version
	if meta.version >= 2.0 {
		meta.numberWidth = 8
		meta.numberFormat = NumfmtBe8bytesq
	} else {
		meta.numberWidth = 4
		meta.numberFormat = NumfmtBe4bytesi
	}

	// Process encoding
	encoding := strings.ToLower(headerInfo.Encoding)
	switch encoding {
	case "gbk", "gb2312":
		meta.encoding = EncodingGb18030
	case "big5":
		meta.encoding = EncodingBig5
	case "utf-16":
		meta.encoding = EncodingUtf16
	default:
		meta.encoding = EncodingUtf8
	}

	// Correct encoding for MDD type
	if mdict.fileType == MdictTypeMdd {
		meta.encoding = EncodingUtf16
	}

	// 4 bytes header length + header byte size + 4 bytes adler checksum
	meta.keyBlockMetaStartOffset = int64(4 + dictHeader.headerBytesSize + 4)

	meta.description = headerInfo.Description
	meta.title = headerInfo.Title
	meta.creationDate = headerInfo.CreationDate
	meta.generatedByEngineVersion = headerInfo.GeneratedByEngineVersion

	mdict.meta = meta

	return nil
}

// readMDictFileHeader reads the raw header data block from an MDict file.
func readMDictFileHeader(filename string) (*mdictHeader, error) {
	log.Debugf("Reading MDict header from file: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", filename, err)
	}
	defer file.Close()

	var dictHeaderPartByteSize int64

	// Read header length
	var headerBytesSize uint32
	dictHeaderPartByteSize += 4
	if err := binary.Read(file, binary.BigEndian, &headerBytesSize); err != nil {
		return nil, fmt.Errorf("failed to read header length from '%s': %w", filename, err)
	}
	log.Debugf("File '%s': Header length: %d", filename, headerBytesSize)

	// Read header info bytes
	headerInfoBytes := make([]byte, headerBytesSize)
	dictHeaderPartByteSize += int64(headerBytesSize)
	if _, err := file.Read(headerInfoBytes); err != nil {
		return nil, fmt.Errorf("failed to read header info bytes from '%s': %w", filename, err)
	}

	// Read adler32 checksum
	var adler32Checksum uint32
	dictHeaderPartByteSize += 4
	if err := binary.Read(file, binary.BigEndian, &adler32Checksum); err != nil {
		return nil, fmt.Errorf("failed to read adler32 checksum from '%s': %w", filename, err)
	}
	log.Debugf("File '%s': Header adler32 checksum from file: %d", filename, adler32Checksum)

	// Convert UTF-16LE encoded header bytes to UTF-8 string
	utfHeaderInfo := littleEndianBinUTF16ToUTF8(headerInfoBytes, 0, int(headerBytesSize))
	// Compatibility fix: replace "Library_Data" with "Dictionary"
	utfHeaderInfo = strings.Replace(utfHeaderInfo, "Library_Data", "Dictionary", 1)

	mdict := &mdictHeader{
		headerBytesSize:          headerBytesSize,
		headerInfoBytes:          headerInfoBytes,
		headerInfo:               utfHeaderInfo,
		adler32Checksum:          adler32Checksum,
		dictionaryHeaderByteSize: dictHeaderPartByteSize,
	}

	return mdict, nil
}

// readKeyBlockMeta reads the metadata of the key block.
func (mdict *MdictBase) readKeyBlockMeta() error {
	log.Infof("Reading key block metadata: %s", mdict.filePath)
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for reading key block metadata: %w", mdict.filePath, err)
	}
	defer file.Close()

	keyBlockMeta := &mdictKeyBlockMeta{}
	log.Debugf("Key block metadata read settings for '%s'. Version: %.1f, NumberWidth: %d, KeyBlockMetaStartOffset: %d",
		mdict.filePath, mdict.meta.version, mdict.meta.numberWidth, mdict.meta.keyBlockMetaStartOffset)

	// Key block metadata section
	// If version > 2.0, the key block metadata section is 40 bytes long
	// Otherwise, it is 16 bytes
	keyBlockMetaBytesNum := 0
	if mdict.meta.version >= 2.0 {
		keyBlockMetaBytesNum = 8 * 5
	} else {
		keyBlockMetaBytesNum = 4 * 4
	}

	// Key block metadata buffer
	keyBlockMetaBuffer, err := readFileFromPos(file, mdict.meta.keyBlockMetaStartOffset, int64(keyBlockMetaBytesNum))
	if err != nil {
		return fmt.Errorf("failed to read key block metadata buffer for '%s': %w", mdict.filePath, err)
	}

	// If encryption type is EncryptKeyInfoEnc, decrypt the key block metadata
	if mdict.meta.encryptType == EncryptKeyInfoEnc {
		log.Debugf("Key block metadata is encrypted (type %d) for '%s', decrypting...", mdict.meta.encryptType, mdict.filePath)
		decryptedKeyBlockMetaBuffer := mdxDecrypt(keyBlockMetaBuffer, int64(keyBlockMetaBytesNum))
		if len(decryptedKeyBlockMetaBuffer) != keyBlockMetaBytesNum {
			return fmt.Errorf("key block metadata decryption error for '%s': output size mismatch (expected %d, got %d)", mdict.filePath, keyBlockMetaBytesNum, len(decryptedKeyBlockMetaBuffer))
		}
		keyBlockMetaBuffer = decryptedKeyBlockMetaBuffer
		log.Debugf("Key block metadata decryption complete for: %s", mdict.filePath)
	}

	// 1. [0:8]([0:4]) - Number of key blocks
	keyBlockNumBytes := keyBlockMetaBuffer[0:mdict.meta.numberWidth]

	var keyBlockNumber uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockNumber = beBinToU64(keyBlockNumBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockNumber = uint64(beBinToU32(keyBlockNumBytes))
	}
	keyBlockMeta.keyBlockNum = int64(keyBlockNumber)

	// 2. [8:16]([4:8]) - Number of entries
	entriesNumBytes := keyBlockMetaBuffer[mdict.meta.numberWidth : mdict.meta.numberWidth+mdict.meta.numberWidth]
	var entriesNum uint64
	if mdict.meta.numberWidth == 8 {
		entriesNum = beBinToU64(entriesNumBytes)
	} else if mdict.meta.numberWidth == 4 {
		entriesNum = uint64(beBinToU32(entriesNumBytes))
	}
	keyBlockMeta.entriesNum = int64(entriesNum)

	var keyBlockInfoSizeBytesStartOffset int

	// 3. [16:24] - Decompressed size of key block info (only if version >= 2.0)
	if mdict.meta.version >= 2.0 {
		keyBlockInfoDecompressSizeBytes := keyBlockMetaBuffer[mdict.meta.numberWidth*2 : mdict.meta.numberWidth*2+mdict.meta.numberWidth]

		var keyBlockInfoDecompressSize uint64
		if mdict.meta.numberWidth == 8 {
			keyBlockInfoDecompressSize = beBinToU64(keyBlockInfoDecompressSizeBytes)
		} else if mdict.meta.numberWidth == 4 {
			keyBlockInfoDecompressSize = uint64(beBinToU32(keyBlockInfoDecompressSizeBytes))
		}
		keyBlockMeta.keyBlockInfoDecompressSize = int64(keyBlockInfoDecompressSize)

		keyBlockInfoSizeBytesStartOffset = mdict.meta.numberWidth * 3

	} else {
		keyBlockInfoSizeBytesStartOffset = mdict.meta.numberWidth * 2
	}

	// 4. [24:32]([8:12]) - Size of key block info
	keyBlockInfoSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset : keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth]

	var keyBlockInfoSize uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockInfoSize = beBinToU64(keyBlockInfoSizeBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockInfoSize = uint64(beBinToU32(keyBlockInfoSizeBytes))
	}

	keyBlockMeta.keyBlockInfoCompressedSize = int64(keyBlockInfoSize)

	// 5. [32:40]([12:16]) - Size of key blocks
	keyBlockDataSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth : keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth+mdict.meta.numberWidth]

	var keyBlockDataSize uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockDataSize = beBinToU64(keyBlockDataSizeBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockDataSize = uint64(beBinToU32(keyBlockDataSizeBytes))
	}
	keyBlockMeta.keyBlockDataTotalSize = int64(keyBlockDataSize)

	// 6. [40:44] - 4-byte checksum (TODO: skip if version > 2.0)
	// TODO: checksum verification

	if mdict.meta.version >= 2.0 {
		keyBlockMeta.keyBlockInfoStartOffset = mdict.meta.keyBlockMetaStartOffset + 40 + 4
	} else {
		keyBlockMeta.keyBlockInfoStartOffset = mdict.meta.keyBlockMetaStartOffset + 16
	}

	mdict.keyBlockMeta = keyBlockMeta

	return nil
}

func (mdict *MdictBase) readKeyBlockInfo() error {
	log.Debugf("Reading key block info: %s", mdict.filePath)
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for reading key block info: %w", mdict.filePath, err)
	}
	defer file.Close()

	log.Debugf("Reading key block info from offset %d, size %d for '%s'", mdict.keyBlockMeta.keyBlockInfoStartOffset, mdict.keyBlockMeta.keyBlockInfoCompressedSize, mdict.filePath)
	buffer, err := readFileFromPos(file, mdict.keyBlockMeta.keyBlockInfoStartOffset, mdict.keyBlockMeta.keyBlockInfoCompressedSize)
	if err != nil {
		return fmt.Errorf("failed to read key block info data for '%s': %w", mdict.filePath, err)
	}

	if err := mdict.decodeKeyBlockInfo(buffer); err != nil {
		return fmt.Errorf("failed to decode key block info for '%s': %w", mdict.filePath, err)
	}
	log.Debugf("Key block info successfully read and decoded for: %s", mdict.filePath)
	return nil
}

func (mdict *MdictBase) decodeKeyBlockInfo(data []byte) error {
	log.Debugf("Decoding key block info for '%s'. Data length: %d, EncryptType: %d", mdict.filePath, len(data), mdict.meta.encryptType)
	if len(data) < 8 {
		return fmt.Errorf("key block info data is too short for magic and checksum check (len %d) for '%s'", len(data), mdict.filePath)
	}

	// Decrypt
	var keyBlockInfoDecryptedBuffer []byte
	if mdict.meta.encryptType == EncryptKeyInfoEnc {
		log.Debugf("Key block info for '%s' is encrypted (type %d), decrypting. Compressed size: %d", mdict.filePath, mdict.meta.encryptType, mdict.keyBlockMeta.keyBlockInfoCompressedSize)
		decryptedBuffer := mdxDecrypt(data, mdict.keyBlockMeta.keyBlockInfoCompressedSize)
		if len(decryptedBuffer) != int(mdict.keyBlockMeta.keyBlockInfoCompressedSize) {
			return fmt.Errorf("key block info decryption error for '%s': output size mismatch (expected %d, got %d)", mdict.filePath, mdict.keyBlockMeta.keyBlockInfoCompressedSize, len(decryptedBuffer))
		}
		keyBlockInfoDecryptedBuffer = decryptedBuffer
		log.Debugf("Key block info decryption complete for: %s", mdict.filePath)
	} else {
		keyBlockInfoDecryptedBuffer = data
	}

	// Check compression type
	compressionType := keyBlockInfoDecryptedBuffer[0:4]
	if !(compressionType[0] == 2 && compressionType[1] == 0 && compressionType[2] == 0 && compressionType[3] == 0) {
		log.Warningf("Compression type of key block info for '%s' is not zlib [02000000], but: %x", mdict.filePath, compressionType)
		// Some dictionaries may not have a compression type header and are just data.
	}

	// Decompress ZLIB data
	expectedChecksum := beBinToU32(keyBlockInfoDecryptedBuffer[4:8])
	compressedDataToDecompress := keyBlockInfoDecryptedBuffer[8:]
	decompressedKeyInfoBuffer, err := zlibDecompress(compressedDataToDecompress, 0, int64(len(compressedDataToDecompress)))
	if err != nil {
		return fmt.Errorf("ZLIB decompression of key block info failed: %w", err)
	}

	if int64(len(decompressedKeyInfoBuffer)) != mdict.keyBlockMeta.keyBlockInfoDecompressSize {
		return fmt.Errorf("ZLIB decompressed size of key block info mismatch: expected %d, got %d",
			mdict.keyBlockMeta.keyBlockInfoDecompressSize, len(decompressedKeyInfoBuffer))
	}

	// Calculate and verify Alder32 checksum on the decompressed data
	actualChecksum := adler32.Checksum(decompressedKeyInfoBuffer)
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("key block info checksum mismatch: expected %d, got %d", expectedChecksum, actualChecksum)
	}

	// Decode key block entries
	var counter int64
	var currentEntriesSize int64
	var numEntriesCounter int64
	byteWidth := 1
	textTerm := 0

	if mdict.meta.version >= 2.0 {
		byteWidth = 2
		textTerm = 1
	}

	var dataOffset = 0
	var compressSizeAccumulator = 0
	var decompressSizeAccumulator = 0

	keyBlockInfo := &mdictKeyBlockInfo{
		keyBlockEntriesStartOffset: 0,
		keyBlockInfoList:           make([]*mdictKeyBlockInfoItem, 0),
	}

	for counter < mdict.keyBlockMeta.keyBlockNum {
		var firstKeySize, lastKeySize int
		var firstKey, lastKey string

		if mdict.meta.version >= 2.0 {
			currentEntriesSize = int64(beBinToU64(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
			dataOffset += mdict.meta.numberWidth
			firstKeySize = int(beBinToU16(decompressedKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
			dataOffset += byteWidth
		} else {
			currentEntriesSize = int64(beBinToU32(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
			dataOffset += mdict.meta.numberWidth
			firstKeySize = int(beBinToU8(decompressedKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
			dataOffset += byteWidth
		}
		numEntriesCounter += currentEntriesSize

		var stepGap = 0
		var termSize = textTerm
		if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
			stepGap = (firstKeySize + textTerm) * 2
			termSize = textTerm * 2
		} else {
			stepGap = firstKeySize + textTerm
			termSize = textTerm
		}

		firstKey = bigEndianBinToUTF8(decompressedKeyInfoBuffer, dataOffset, stepGap-termSize)

		dataOffset += stepGap

		if mdict.meta.version >= 2.0 {
			lastKeySize = int(beBinToU16(decompressedKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		} else {
			lastKeySize = int(beBinToU8(decompressedKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		}
		dataOffset += byteWidth

		if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
			stepGap = (lastKeySize + textTerm) * 2
			termSize = textTerm * 2
		} else {
			stepGap = lastKeySize + textTerm
			termSize = textTerm
		}

		lastKey = bigEndianBinToUTF8(decompressedKeyInfoBuffer, dataOffset, stepGap-termSize)

		dataOffset += stepGap
		var keyBlockCompressSize int
		if mdict.meta.version >= 2.0 {
			keyBlockCompressSize = int(beBinToU64(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		} else {
			keyBlockCompressSize = int(beBinToU32(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		}
		dataOffset += mdict.meta.numberWidth

		var keyBlockDecompressSize int
		if mdict.meta.version >= 2.0 {
			keyBlockDecompressSize = int(beBinToU64(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		} else {
			keyBlockDecompressSize = int(beBinToU32(decompressedKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		}

		dataOffset += mdict.meta.numberWidth

		keyBlockInfoItem := &mdictKeyBlockInfoItem{
			firstKey:                      firstKey,
			firstKeySize:                  firstKeySize,
			lastKey:                       lastKey,
			lastKeySize:                   lastKeySize,
			keyBlockInfoIndex:             int(counter),
			keyBlockCompressSize:          int64(keyBlockCompressSize),
			keyBlockCompAccumulator:       int64(compressSizeAccumulator),
			keyBlockDeCompressSize:        int64(keyBlockDecompressSize),
			keyBlockDeCompressAccumulator: int64(decompressSizeAccumulator),
		}

		compressSizeAccumulator += keyBlockCompressSize
		decompressSizeAccumulator += keyBlockDecompressSize

		keyBlockInfo.keyBlockInfoList = append(keyBlockInfo.keyBlockInfoList, keyBlockInfoItem)

		counter++

	}
	keyBlockInfo.keyBlockEntriesStartOffset = mdict.keyBlockMeta.keyBlockInfoCompressedSize + mdict.keyBlockMeta.keyBlockInfoStartOffset

	mdict.keyBlockInfo = keyBlockInfo

	if int64(compressSizeAccumulator) != mdict.keyBlockMeta.keyBlockDataTotalSize {
		return fmt.Errorf("key block data compressed size mismatch with metadata (%d/%d)", compressSizeAccumulator, mdict.keyBlockMeta.keyBlockDataTotalSize)
	}

	return nil

}

func (mdict *MdictBase) readKeyEntries() error {
	log.Debugf("Reading key entries: %s", mdict.filePath)
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for reading key entries: %w", mdict.filePath, err)
	}
	defer file.Close()

	log.Debugf("Reading key entries data for '%s' from offset %d, total size %d", mdict.filePath, mdict.keyBlockInfo.keyBlockEntriesStartOffset, mdict.keyBlockMeta.keyBlockDataTotalSize)
	buffer, err := readFileFromPos(file,
		mdict.keyBlockInfo.keyBlockEntriesStartOffset,
		mdict.keyBlockMeta.keyBlockDataTotalSize)
	if err != nil {
		return fmt.Errorf("failed to read key entries data for '%s': %w", mdict.filePath, err)
	}

	if err := mdict.decodeKeyEntries(buffer); err != nil {
		return fmt.Errorf("failed to decode key entries for '%s': %w", mdict.filePath, err)
	}
	log.Debugf("Key entries successfully read and decoded for '%s'. Total entries: %d", mdict.filePath, mdict.keyBlockData.keyEntriesSize)
	return nil
}

func (mdict *MdictBase) decodeKeyEntries(keyBlockDataCompressBuffer []byte) error {
	log.Debugf("Decoding key entries for '%s'. Total compressed data length: %d", mdict.filePath, len(keyBlockDataCompressBuffer))
	var start, end, compAccu int64

	keyBlockData := &mdictKeyBlockData{
		keyEntries:                 make([]*MDictKeywordEntry, 0),
		keyEntriesSize:             0,
		recordBlockMetaStartOffset: 0,
	}

	for idx, infoItem := range mdict.keyBlockInfo.keyBlockInfoList {
		compressedSize := infoItem.keyBlockCompressSize
		decompressedSize := infoItem.keyBlockDeCompressSize

		compAccu += infoItem.keyBlockCompressSize
		end = start + compressedSize

		if start != infoItem.keyBlockCompAccumulator {
			return fmt.Errorf("[%d] the key-block data start offset not equal to key block compress accumulator(%d/%d/%d)",
				idx, start, infoItem.keyBlockCompAccumulator, compAccu)
		}

		kbCompType := keyBlockDataCompressBuffer[start : start+4]
		expectedKeyBlockChecksum := beBinToU32(keyBlockDataCompressBuffer[start+4 : start+8])

		var keyBlock []byte

		switch kbCompType[0] {
		case 0: // No compression
			keyBlock = keyBlockDataCompressBuffer[start+8 : end]
		case 1: // LZO compression
			compressedLZOData := keyBlockDataCompressBuffer[start+8 : end]
			reader := bytes.NewReader(compressedLZOData)
			out, err1 := lzo.Decompress1X(reader, 0, int(decompressedSize))
			if err1 != nil {
				return fmt.Errorf("LZO decompression failed for key block %d: %w", idx, err1)
			}
			if len(out) != int(decompressedSize) {
				return fmt.Errorf("LZO decompression output size mismatch for key block %d: expected %d, got %d", idx, decompressedSize, len(out))
			}
			keyBlock = out
		case 2: // ZLIB compression
			compressedZLIBData := keyBlockDataCompressBuffer[start+8 : end]
			out, err2 := zlibDecompress(compressedZLIBData, 0, int64(len(compressedZLIBData)))
			if err2 != nil {
				return err2
			}
			keyBlock = out

			actualKeyBlockChecksum := adler32.Checksum(keyBlock)
			if actualKeyBlockChecksum != expectedKeyBlockChecksum {
				return fmt.Errorf("key block data checksum mismatch for block %d: expected %d, got %d", idx, expectedKeyBlockChecksum, actualKeyBlockChecksum)
			}
		default:
			return fmt.Errorf("cannot determine the compress type %v", kbCompType)
		}

		splitKeys := mdict.splitKeyBlock(keyBlock)

		keyBlockData.keyEntries = append(keyBlockData.keyEntries, splitKeys...)
		keyBlockData.keyEntriesSize += int64(len(splitKeys))

		start = end
	}

	if keyBlockData.keyEntriesSize != mdict.keyBlockMeta.entriesNum {
		return fmt.Errorf("decoded key list items count %d not equal to expected entries number %d for '%s'", keyBlockData.keyEntriesSize, mdict.keyBlockMeta.entriesNum, mdict.filePath)
	}
	keyBlockData.recordBlockMetaStartOffset = mdict.keyBlockInfo.keyBlockEntriesStartOffset + mdict.keyBlockMeta.keyBlockDataTotalSize

	mdict.keyBlockData = keyBlockData

	return nil
}

func (mdict *MdictBase) splitKeyBlock(keyBlock []byte) []*MDictKeywordEntry {
	width := 1
	if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
		width = 2
	}

	var keyList []*MDictKeywordEntry
	keyStartIndex := 0

	for keyStartIndex < len(keyBlock) {
		var recordStartOffset int64
		if mdict.meta.numberWidth == 8 {
			recordStartOffset = int64(beBinToU64(keyBlock[keyStartIndex : keyStartIndex+mdict.meta.numberWidth]))
		} else {
			recordStartOffset = int64(beBinToU32(keyBlock[keyStartIndex : keyStartIndex+mdict.meta.numberWidth]))
		}

		keyEndIndex := keyStartIndex + mdict.meta.numberWidth
		for i := keyEndIndex; i < len(keyBlock); i += width {
			if (width == 1 && keyBlock[i] == 0) || (width == 2 && keyBlock[i] == 0 && keyBlock[i+1] == 0) {
				keyEndIndex = i
				break
			}
		}

		keyTextBytes := keyBlock[keyStartIndex+mdict.meta.numberWidth : keyEndIndex]
		keyText := string(keyTextBytes)
		var err error

		if mdict.meta.encoding == EncodingUtf16 {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				keyText = string(keyTextBytes)
			}
		}

		if mdict.fileType == MdictTypeMdd {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				log.Errorf("Error decoding UTF-16 for MDD key text (offset %d): %v. KeyTextBytes: %x", keyStartIndex+mdict.meta.numberWidth, err, keyTextBytes)
				keyText = string(keyTextBytes) // Fallback to raw string to avoid panic
			}
		}

		keyStartIndex = keyEndIndex + width
		entry := &MDictKeywordEntry{
			RecordStartOffset: recordStartOffset,
			KeyWord:           keyText,
			KeyBlockIdx:       int64(keyStartIndex),
		}
		if len(keyList) > 0 {
			keyList[len(keyList)-1].RecordEndOffset = entry.RecordStartOffset
		}
		keyList = append(keyList, entry)
	}

	return keyList
}

func (mdict *MdictBase) readRecordBlockMeta() error {
	log.Debugf("Reading record block metadata for: %s", mdict.filePath)
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for record block metadata: %w", mdict.filePath, err)
	}
	defer file.Close()

	recordBlockMetaBufferLen := int64(16)
	if mdict.meta.version >= 2.0 {
		recordBlockMetaBufferLen = 32
	}

	recordBlockStartOffset := mdict.keyBlockInfo.keyBlockEntriesStartOffset + mdict.keyBlockMeta.keyBlockDataTotalSize
	log.Debugf("Reading record block metadata for '%s' from offset %d, length %d", mdict.filePath, recordBlockStartOffset, recordBlockMetaBufferLen)

	buffer, err := readFileFromPos(file, recordBlockStartOffset, recordBlockMetaBufferLen)
	if err != nil {
		return fmt.Errorf("failed to read record block metadata for '%s': %w", mdict.filePath, err)
	}

	if err := mdict.decodeRecordBlockMeta(buffer, recordBlockStartOffset, recordBlockStartOffset+recordBlockMetaBufferLen); err != nil {
		return fmt.Errorf("failed to decode record block metadata for '%s': %w", mdict.filePath, err)
	}
	log.Debugf("Record block metadata successfully read and decoded for '%s': %+v", mdict.filePath, mdict.recordBlockMeta)
	return nil
}

func (mdict *MdictBase) decodeRecordBlockMeta(data []byte, startOffset, endOffset int64) error {
	recordBlockMeta := &mdictRecordBlockMeta{
		keyRecordMetaStartOffset: startOffset,
		keyRecordMetaEndOffset:   endOffset,
	}

	offset := 0

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockNum = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockNum = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
	}
	offset += mdict.meta.numberWidth

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.entriesNum = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.entriesNum = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
	}
	if recordBlockMeta.entriesNum != mdict.keyBlockMeta.entriesNum {
		return fmt.Errorf("record block entries number %d does not match key block entries number %d for '%s'", recordBlockMeta.entriesNum, mdict.keyBlockMeta.entriesNum, mdict.filePath)
	}
	offset += mdict.meta.numberWidth

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockInfoCompSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockInfoCompSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
	}
	offset += mdict.meta.numberWidth

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockCompSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockCompSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
	}

	mdict.recordBlockMeta = recordBlockMeta
	return nil
}

func (mdict *MdictBase) readRecordBlockInfo() error {
	log.Debugf("Reading record block info for: %s", mdict.filePath)
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for record block info: %w", mdict.filePath, err)
	}
	defer file.Close()

	recordBlockInfoStartOffset := mdict.recordBlockMeta.keyRecordMetaEndOffset
	recordBlockInfoLen := mdict.recordBlockMeta.recordBlockInfoCompSize
	log.Debugf("Reading record block info for '%s' from offset %d, length %d", mdict.filePath, recordBlockInfoStartOffset, recordBlockInfoLen)

	buffer, err := readFileFromPos(file, recordBlockInfoStartOffset, recordBlockInfoLen)
	if err != nil {
		return fmt.Errorf("failed to read record block info data for '%s': %w", mdict.filePath, err)
	}

	if err := mdict.decodeRecordBlockInfo(buffer, recordBlockInfoStartOffset, recordBlockInfoStartOffset+recordBlockInfoLen); err != nil {
		return fmt.Errorf("failed to decode record block info for '%s': %w", mdict.filePath, err)
	}
	log.Debugf("Record block info successfully read and decoded for '%s'. Number of record blocks: %d", mdict.filePath, len(mdict.recordBlockInfo.recordInfoList))
	return nil
}

func (mdict *MdictBase) decodeRecordBlockInfo(data []byte, startOffset, endOffset int64) error {
	var recordBlockInfoList []*MdictRecordBlockInfoListItem
	var offset int
	var compAccu, decompAccu int64

	for i := int64(0); i < mdict.recordBlockMeta.recordBlockNum; i++ {
		var compSize, decompSize int64
		if mdict.meta.version >= 2.0 {
			compSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
		} else {
			compSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
		}
		offset += mdict.meta.numberWidth

		if mdict.meta.version >= 2.0 {
			decompSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
		} else {
			decompSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
		}
		offset += mdict.meta.numberWidth

		recordBlockInfoList = append(recordBlockInfoList, &MdictRecordBlockInfoListItem{
			compressSize:                compSize,
			deCompressSize:              decompSize,
			compressAccumulatorOffset:   compAccu,
			deCompressAccumulatorOffset: decompAccu,
		})

		compAccu += compSize
		decompAccu += decompSize
	}

	if int64(len(recordBlockInfoList)) != mdict.recordBlockMeta.recordBlockNum {
		return fmt.Errorf("decoded record block info items count %d not equal to expected record block number %d for '%s'. CompAccumulator: %d, DecompAccumulator: %d",
			len(recordBlockInfoList), mdict.recordBlockMeta.recordBlockNum, mdict.filePath, compAccu, decompAccu)
	}
	if int64(offset) != mdict.recordBlockMeta.recordBlockInfoCompSize {
		return fmt.Errorf("record block info decoded offset %d not equal to expected compressed size %d for '%s'", offset, mdict.recordBlockMeta.recordBlockInfoCompSize, mdict.filePath)
	}
	if compAccu != mdict.recordBlockMeta.recordBlockCompSize {
		return fmt.Errorf("record block info accumulated compressed size %d not equal to expected total compressed size %d for '%s'", compAccu, mdict.recordBlockMeta.recordBlockCompSize, mdict.filePath)
	}

	mdict.recordBlockInfo = &mdictRecordBlockInfo{
		recordInfoList:             recordBlockInfoList,
		recordBlockInfoStartOffset: startOffset,
		recordBlockInfoEndOffset:   endOffset,
		recordBlockDataStartOffset: endOffset,
	}

	return nil
}

func (mdict *MdictBase) buildRecordRangeTree() {
	log.Debugf("Building record range tree with %d items for: %s", len(mdict.recordBlockInfo.recordInfoList), mdict.filePath)
	BuildRangeTree(mdict.recordBlockInfo.recordInfoList, mdict.rangeTreeRoot)
	log.Debugf("Record range tree built for: %s", mdict.filePath)
}

func (mdict *MdictBase) keywordEntryToIndex(item *MDictKeywordEntry) (*MDictKeywordIndex, error) {
	var recordBlockInfo *MdictRecordBlockInfoListItem

	if mdict.rangeTreeRoot != nil {
		log.Debugf("Attempting to find record block info for offset %d using range tree.", item.RecordStartOffset)
		rbInfo := QueryRangeData(mdict.rangeTreeRoot, item.RecordStartOffset)
		if rbInfo != nil {
			recordBlockInfo = rbInfo
			log.Debugf("Found record block info for offset %d using range tree: %+v", item.RecordStartOffset, recordBlockInfo)
		} else {
			log.Debugf("Record block info for offset %d not found using range tree. Will attempt linear scan.", item.RecordStartOffset)
		}
	} else {
		log.Debugf("Range tree not initialized. Using linear scan for offset %d.", item.RecordStartOffset)
	}

	if recordBlockInfo == nil {
		log.Debugf("Performing linear scan for record block info for offset %d.", item.RecordStartOffset)
		var found bool
		for i, rbi := range mdict.recordBlockInfo.recordInfoList {
			if item.RecordStartOffset >= rbi.deCompressAccumulatorOffset && item.RecordStartOffset < (rbi.deCompressAccumulatorOffset+rbi.deCompressSize) {
				recordBlockInfo = rbi
				log.Debugf("Found record block info via linear scan at index %d for offset %d: %+v", i, item.RecordStartOffset, recordBlockInfo)
				found = true
				break
			}
		}
		if !found {
			log.Errorf("Linear scan failed to find record block info for offset %d for '%s'. Total record blocks: %d.", item.RecordStartOffset, mdict.filePath, len(mdict.recordBlockInfo.recordInfoList))
			if len(mdict.recordBlockInfo.recordInfoList) > 0 {
				log.Debugf("First record block info for linear scan failure (file '%s'): %+v", mdict.filePath, mdict.recordBlockInfo.recordInfoList[0])
				log.Debugf("Last record block info for linear scan failure (file '%s'): %+v", mdict.filePath, mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1])
			}
			return nil, fmt.Errorf("key-item's record block info not found for RecordStartOffset %d via linear scan for '%s'", item.RecordStartOffset, mdict.filePath)
		}
	}

	recordBlockFileOffset := recordBlockInfo.compressAccumulatorOffset + mdict.recordBlockInfo.recordBlockDataStartOffset
	keywordStartOffsetInDecompressedBlock := item.RecordStartOffset - recordBlockInfo.deCompressAccumulatorOffset
	var keywordEndOffsetInDecompressedBlock int64
	if item.RecordEndOffset == 0 {
		keywordEndOffsetInDecompressedBlock = recordBlockInfo.deCompressSize
		log.Debugf("RecordEndOffset is 0, setting keyword end to deCompressSize: %d for item offset %d", keywordEndOffsetInDecompressedBlock, item.RecordStartOffset)
	} else {
		keywordEndOffsetInDecompressedBlock = item.RecordEndOffset - recordBlockInfo.deCompressAccumulatorOffset
		log.Debugf("RecordEndOffset is %d, calculated keyword end to %d for item offset %d", item.RecordEndOffset, keywordEndOffsetInDecompressedBlock, item.RecordStartOffset)
	}

	if keywordStartOffsetInDecompressedBlock < 0 || keywordStartOffsetInDecompressedBlock > recordBlockInfo.deCompressSize {
		log.Errorf("Calculated keyword start offset %d is out of bounds for its decompressed record block (size %d). Item: %+v, RecordBlockInfo: %+v",
			keywordStartOffsetInDecompressedBlock, recordBlockInfo.deCompressSize, item, recordBlockInfo)
		return nil, fmt.Errorf("calculated keyword start offset %d is out of bounds for its decompressed record block (size %d, item offset %d, block decompress acc offset %d)",
			keywordStartOffsetInDecompressedBlock, recordBlockInfo.deCompressSize, item.RecordStartOffset, recordBlockInfo.deCompressAccumulatorOffset)
	}
	if keywordEndOffsetInDecompressedBlock < keywordStartOffsetInDecompressedBlock || keywordEndOffsetInDecompressedBlock > recordBlockInfo.deCompressSize {
		log.Errorf("Calculated keyword end offset %d is out of bounds (start: %d, decompressed size: %d). Item: %+v, RecordBlockInfo: %+v",
			keywordEndOffsetInDecompressedBlock, keywordStartOffsetInDecompressedBlock, recordBlockInfo.deCompressSize, item, recordBlockInfo)
		return nil, fmt.Errorf("calculated keyword end offset %d is out of bounds (start: %d, decompressed size: %d, item end offset %d, block decompress acc offset %d)",
			keywordEndOffsetInDecompressedBlock, keywordStartOffsetInDecompressedBlock, recordBlockInfo.deCompressSize, item.RecordEndOffset, recordBlockInfo.deCompressAccumulatorOffset)
	}

	log.Debugf("Successfully created MDictKeywordIndex for item offset %d: StartInFile=%d, KeyWordStartOffset=%d, KeyWordEndOffset=%d",
		item.RecordStartOffset, recordBlockFileOffset, keywordStartOffsetInDecompressedBlock, keywordEndOffsetInDecompressedBlock)

	return &MDictKeywordIndex{
		KeywordEntry: *item,
		RecordBlock: MDictKeywordIndexRecordBlock{
			DataStartOffset:          recordBlockFileOffset,
			CompressSize:             recordBlockInfo.compressSize,
			DeCompressSize:           recordBlockInfo.deCompressSize,
			KeyWordPartStartOffset:   keywordStartOffsetInDecompressedBlock,
			KeyWordPartDataEndOffset: keywordEndOffsetInDecompressedBlock,
		},
	}, nil
}

func (mdict *MdictBase) locateByKeywordIndex(index *MDictKeywordIndex) ([]byte, error) {
	return _locateDefByKWIndexInternal(index,
		mdict.filePath,
		mdict.meta.encryptType == EncryptRecordEnc,
		mdict.fileType == MdictTypeMdd,
		mdict.meta.encoding == EncodingUtf16,
		index.KeywordEntry.KeyWord,
	)
}

func _fetchAndDecodeRecordBlock(filePath string, fileOffset int64, compressedSize int64, decompressedSize int64, isEncrypted bool, keywordForLog string) ([]byte, error) {
	log.Debugf("Fetching record block for keyword '%s': fileOffset=%d, compressedSize=%d, decompressedSize=%d, isEncrypted=%v",
		keywordForLog, fileOffset, compressedSize, decompressedSize, isEncrypted)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", filePath, err)
	}
	defer file.Close()

	recordBlockDataCompBuff, err := readFileFromPos(file, fileOffset, compressedSize)
	if err != nil {
		return nil, fmt.Errorf("error reading record block data for keyword '%s' from offset %d, size %d: %w", keywordForLog, fileOffset, compressedSize, err)
	}
	if recordBlockDataCompBuff == nil {
		return nil, fmt.Errorf("read empty record block data for keyword '%s' from offset %d, size %d", keywordForLog, fileOffset, compressedSize)
	}
	if len(recordBlockDataCompBuff) < 8 {
		return nil, fmt.Errorf("record block data for keyword '%s' is too short (%d bytes) to contain header", keywordForLog, len(recordBlockDataCompBuff))
	}

	rbCompType := recordBlockDataCompBuff[0:4]
	expectedChecksum := beBinToU32(recordBlockDataCompBuff[4:8])
	log.Debugf("Record block for '%s': CompType=%v, ExpectedChecksum=%d", keywordForLog, rbCompType, expectedChecksum)

	var recordBlock []byte
	var dataToProcess []byte

	if isEncrypted {
		log.Debugf("Decrypting record block for '%s'", keywordForLog)
		decryptedFullBlock := mdxDecrypt(recordBlockDataCompBuff, compressedSize)
		if int64(len(decryptedFullBlock)) != compressedSize {
			return nil, fmt.Errorf("decryption error for record block of '%s': output size mismatch, expected %d, got %d", keywordForLog, compressedSize, len(decryptedFullBlock))
		}
		if len(decryptedFullBlock) < 8 {
			return nil, fmt.Errorf("decrypted record block for '%s' is too short (%d bytes)", keywordForLog, len(decryptedFullBlock))
		}
		dataToProcess = decryptedFullBlock[8:]
	} else {
		dataToProcess = recordBlockDataCompBuff[8:]
	}

	switch rbCompType[0] {
	case 0: // No compression
		log.Debugf("Record block for '%s' is not compressed.", keywordForLog)
		recordBlock = dataToProcess
	case 1: // LZO compression
		log.Debugf("Decompressing LZO record block for '%s'. DataToProcess len: %d, ExpectedDecompressed: %d", keywordForLog, len(dataToProcess), decompressedSize)
		reader := bytes.NewReader(dataToProcess)
		out, err1 := lzo.Decompress1X(reader, 0, int(decompressedSize))
		if err1 != nil {
			return nil, fmt.Errorf("LZO decompression failed for record block of '%s': %w", keywordForLog, err1)
		}
		recordBlock = out
	case 2: // ZLIB compression
		log.Debugf("Decompressing ZLIB record block for '%s'. DataToProcess len: %d, ExpectedDecompressed: %d", keywordForLog, len(dataToProcess), decompressedSize)
		out, err2 := zlibDecompress(dataToProcess, 0, int64(len(dataToProcess)))
		if err2 != nil {
			return nil, fmt.Errorf("ZLIB decompression failed for record block of '%s': %w", keywordForLog, err2)
		}
		recordBlock = out
	default:
		return nil, fmt.Errorf("unknown record block compression type %v for keyword '%s'", rbCompType, keywordForLog)
	}

	actualChecksum := adler32.Checksum(recordBlock)
	if actualChecksum != expectedChecksum {
		log.Errorf("Checksum mismatch for record block of '%s': expected %d, got %d. Decompressed len: %d", keywordForLog, expectedChecksum, actualChecksum, len(recordBlock))
		return nil, fmt.Errorf("record block checksum mismatch for keyword '%s': expected %d, got %d", keywordForLog, expectedChecksum, actualChecksum)
	}
	log.Debugf("Checksum verified for record block of '%s': %d", keywordForLog, actualChecksum)

	if int64(len(recordBlock)) != decompressedSize {
		log.Errorf("Decompressed size mismatch for record block of '%s': expected %d, got %d", keywordForLog, decompressedSize, len(recordBlock))
		return nil, fmt.Errorf("record block decompressed size mismatch for keyword '%s': expected %d, got %d", keywordForLog, decompressedSize, len(recordBlock))
	}
	log.Debugf("Decompressed size verified for record block of '%s': %d", keywordForLog, len(recordBlock))

	return recordBlock, nil
}

func _locateDefByKWIndexInternal(index *MDictKeywordIndex, filePath string, isRecordEncrypted, isMdd, isUtf16 bool, keywordForLog string) ([]byte, error) {
	log.Infof("Locating definition for keyword '%s' using index: %+v", keywordForLog, index.RecordBlock)

	decompressedRecordBlock, err := _fetchAndDecodeRecordBlock(filePath,
		index.RecordBlock.DataStartOffset,
		index.RecordBlock.CompressSize,
		index.RecordBlock.DeCompressSize,
		isRecordEncrypted,
		keywordForLog,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch/decode record block for keyword '%s': %w", keywordForLog, err)
	}

	start := index.RecordBlock.KeyWordPartStartOffset
	end := index.RecordBlock.KeyWordPartDataEndOffset

	if start < 0 || end < start || end > int64(len(decompressedRecordBlock)) {
		log.Errorf("Invalid keyword part offsets for '%s': start=%d, end=%d, block_len=%d. Index: %+v",
			keywordForLog, start, end, len(decompressedRecordBlock), index)
		return nil, fmt.Errorf("invalid keyword part offsets for '%s': start=%d, end=%d, decompressed block length=%d",
			keywordForLog, start, end, len(decompressedRecordBlock))
	}
	log.Debugf("Extracting keyword data for '%s' from decompressed block: start=%d, end=%d", keywordForLog, start, end)
	data := decompressedRecordBlock[start:end]

	if isMdd {
		log.Infof("Returning raw MDD data for '%s' (length %d)", keywordForLog, len(data))
		return data, nil
	}

	if isUtf16 {
		log.Infof("Decoding UTF-16 data for keyword '%s' (original length %d)", keywordForLog, len(data))
		datastr, err1 := decodeLittleEndianUtf16(data)
		if err1 != nil {
			log.Errorf("UTF-16 decoding failed for keyword '%s': %v", keywordForLog, err1)
			return nil, fmt.Errorf("UTF-16 decoding failed for '%s': %w", keywordForLog, err1)
		}
		log.Debugf("Decoded UTF-16 data for '%s' to length %d", keywordForLog, len(datastr))
		return []byte(datastr), nil
	}

	log.Infof("Returning data for keyword '%s' (length %d, encoding assumed UTF-8 or similar)", keywordForLog, len(data))
	return data, nil
}

func locateDefByKWIndex(index *MDictKeywordIndex, filePath string, isRecordEncrypted, isMdd, isUtf16 bool) ([]byte, error) {
	return _locateDefByKWIndexInternal(index, filePath, isRecordEncrypted, isMdd, isUtf16, index.KeywordEntry.KeyWord)
}

func (mdict *MdictBase) locateByKeywordEntry(item *MDictKeywordEntry) ([]byte, error) {
	log.Debugf("Locating by keyword entry: %s (Offset: %d)", item.KeyWord, item.RecordStartOffset)
	index, err := mdict.keywordEntryToIndex(item)
	if err != nil {
		log.Errorf("Failed to get keyword index for entry '%s' (offset %d): %v", item.KeyWord, item.RecordStartOffset, err)
		return nil, fmt.Errorf("failed to get keyword index for '%s': %w", item.KeyWord, err)
	}

	log.Debugf("Obtained keyword index for '%s': %+v", item.KeyWord, index.RecordBlock)

	return _locateDefByKWIndexInternal(index,
		mdict.filePath,
		mdict.meta.encryptType == EncryptRecordEnc,
		mdict.fileType == MdictTypeMdd,
		mdict.meta.encoding == EncodingUtf16,
		item.KeyWord,
	)
}

// GetKeyWordEntries returns all keyword entries in the dictionary.
func (mdict *MdictBase) GetKeyWordEntries() ([]*MDictKeywordEntry, error) {
	return mdict.keyBlockData.keyEntries, nil
}
