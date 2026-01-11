package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf16"
)

var (
	ErrInvalidLSA     = errors.New("invalid LSA file format")
	ErrLSAEntryNotFound = errors.New("LSA entry not found")
)

// LSAEntry LSA 文件中的条目信息
type LSAEntry struct {
	Offset int64
	Size   int64
}

// LSAReader LSA 文件读取器
type LSAReader struct {
	file     *os.File
	index    map[string]LSAEntry // 单词名 -> 条目
	mu       sync.RWMutex
	dataStart int64 // 数据区起始位置
}

// NewLSAReader 创建 LSA 读取器
func NewLSAReader(path string) (*LSAReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := &LSAReader{
		file:  file,
		index: make(map[string]LSAEntry),
	}

	if err := reader.buildIndex(); err != nil {
		file.Close()
		return nil, err
	}

	return reader, nil
}

// buildIndex 解析 LSA 文件并构建索引
// LSA 文件格式（从十六进制分析）：
// - 文件头: "L9SA" (UTF-16LE, 8 bytes)
// - 索引区: 每个条目包含文件名和偏移/大小信息
// - 数据区: WAV 音频数据
func (r *LSAReader) buildIndex() error {
	stat, err := r.file.Stat()
	if err != nil {
		return err
	}
	fileSize := stat.Size()

	// 只读取文件头和索引区（假设索引区不超过前 2MB）
	indexSize := int64(2 * 1024 * 1024)
	if indexSize > fileSize {
		indexSize = fileSize
	}

	data := make([]byte, indexSize)
	if _, err := r.file.ReadAt(data, 0); err != nil {
		return err
	}

	// 验证文件头 "L9SA" (UTF-16LE)
	if len(data) < 8 {
		return ErrInvalidLSA
	}
	
	header := decodeUTF16LE(data[:8])
	if header != "L9SA" {
		return ErrInvalidLSA
	}

	// 使用更高效的方式解析索引
	// 查找所有 ".wav" 的 UTF-16LE 编码位置
	wavPattern := []byte{0x2e, 0x00, 0x77, 0x00, 0x61, 0x00, 0x76, 0x00} // ".wav"
	
	pos := 8
	maxEntries := 50000 // 限制最大条目数防止无限循环
	entryCount := 0

	for pos < len(data)-20 && entryCount < maxEntries {
		// 查找 .wav 模式
		idx := bytes.Index(data[pos:], wavPattern)
		if idx == -1 {
			break
		}
		
		wavPos := pos + idx
		
		// 向前查找文件名开始位置（查找第一个非 UTF-16 有效字符或分隔符）
		nameStart := wavPos
		for i := wavPos - 2; i >= pos && i >= wavPos-200; i -= 2 {
			if i+1 >= len(data) {
				break
			}
			ch := data[i]
			// 检查是否是有效的文件名字符
			isValidChar := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') || ch == ' ' || ch == '-' || ch == '_' || ch == '\''
			if data[i+1] == 0x00 && isValidChar {
				nameStart = i
			} else {
				break
			}
		}
		
		if nameStart >= wavPos {
			pos = wavPos + 8
			continue
		}
		
		// 提取文件名
		nameEnd := wavPos + 8 // 包含 .wav
		if nameEnd > len(data) {
			break
		}
		
		filename := decodeUTF16LE(data[nameStart:nameEnd])
		word := strings.TrimSuffix(filename, ".wav")
		word = strings.ToLower(strings.TrimSpace(word))
		
		if word == "" {
			pos = wavPos + 8
			continue
		}
		
		// 跳过 .wav\r\n\0 寻找偏移量和大小
		searchPos := nameEnd
		if searchPos+12 < len(data) {
			// 查找 0xff 标记后的偏移和大小
			for i := searchPos; i < searchPos+20 && i+10 < len(data); i++ {
				if data[i] == 0xff && i+5 < len(data) {
					offset := int64(binary.LittleEndian.Uint32(data[i+1 : i+5]))
					// 查找下一个 0xff 获取大小
					for j := i + 5; j < i+15 && j+5 < len(data); j++ {
						if data[j] == 0xff {
							size := int64(binary.LittleEndian.Uint32(data[j+1 : j+5]))
							if offset > 0 && size > 0 && offset < fileSize && offset+size <= fileSize {
								r.index[word] = LSAEntry{
									Offset: offset,
									Size:   size,
								}
								entryCount++
							}
							break
						}
					}
					break
				}
			}
		}
		
		pos = wavPos + 8
	}

	return nil
}

// decodeUTF16LE 解码 UTF-16LE 字节为字符串
func decodeUTF16LE(data []byte) string {
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	
	u16s := make([]uint16, len(data)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = binary.LittleEndian.Uint16(data[i*2:])
	}
	
	return string(utf16.Decode(u16s))
}

// Read 读取指定单词的音频数据
func (r *LSAReader) Read(word string) ([]byte, error) {
	word = strings.ToLower(strings.TrimSpace(word))
	
	r.mu.RLock()
	entry, ok := r.index[word]
	r.mu.RUnlock()
	
	if !ok {
		return nil, ErrLSAEntryNotFound
	}
	
	data := make([]byte, entry.Size)
	_, err := r.file.ReadAt(data, entry.Offset)
	if err != nil {
		return nil, err
	}
	
	// 验证是否是有效的 WAV 数据（以 RIFF 开头）
	if len(data) < 4 || !bytes.HasPrefix(data, []byte("RIFF")) {
		return nil, ErrLSAEntryNotFound
	}
	
	return data, nil
}

// ListWords 列出所有可用的单词
func (r *LSAReader) ListWords() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	words := make([]string, 0, len(r.index))
	for word := range r.index {
		words = append(words, word)
	}
	return words
}

// Count 返回索引中的条目数量
func (r *LSAReader) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.index)
}

// Close 关闭文件
func (r *LSAReader) Close() error {
	return r.file.Close()
}

// GetLSAPath 获取 LSA 文件的默认路径
func GetLSAPath(soundDir string) string {
	return filepath.Join(soundDir, "SoundEn.lsa")
}
