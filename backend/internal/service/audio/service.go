package audio

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"dict-hub/internal/service/mdx"
)

var (
	ErrAudioNotFound = errors.New("audio not found")
)

// AudioService 统一音频服务
type AudioService struct {
	mdxManager    mdx.DictManager
	soundDir      string
	lsaReader     *LSAReader
	lsaPath       string
	lsaLoadOnce   sync.Once
	lsaLoadErr    error
	mu            sync.RWMutex
}

// NewAudioService 创建音频服务
func NewAudioService(mdxManager mdx.DictManager, soundDir string) *AudioService {
	svc := &AudioService{
		mdxManager: mdxManager,
		soundDir:   soundDir,
	}

	// 记录 LSA 文件路径，但不在启动时加载（懒加载）
	lsaPath := filepath.Join(soundDir, "SoundEn.lsa")
	if _, err := os.Stat(lsaPath); err == nil {
		svc.lsaPath = lsaPath
		log.Printf("Found LSA file at %s (will load on first access)", lsaPath)
	}

	return svc
}

// ensureLSALoaded 确保 LSA 文件已加载（懒加载）
func (s *AudioService) ensureLSALoaded() (*LSAReader, error) {
	s.lsaLoadOnce.Do(func() {
		if s.lsaPath == "" {
			s.lsaLoadErr = ErrAudioNotFound
			return
		}
		log.Printf("Loading LSA file: %s", s.lsaPath)
		reader, err := NewLSAReader(s.lsaPath)
		if err != nil {
			log.Printf("Failed to load LSA file: %v", err)
			s.lsaLoadErr = err
			return
		}
		s.mu.Lock()
		s.lsaReader = reader
		s.mu.Unlock()
		log.Printf("LSA file loaded successfully with %d entries", reader.Count())
	})
	
	s.mu.RLock()
	reader := s.lsaReader
	s.mu.RUnlock()
	
	if reader == nil {
		return nil, s.lsaLoadErr
	}
	return reader, nil
}

// GetAudio 按优先级查找音频
// 优先级：1. MDD 字典内音频 -> 2. 独立 wav 文件 -> 3. LSA 音频包
func (s *AudioService) GetAudio(word string) (io.Reader, string, error) {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return nil, "", ErrAudioNotFound
	}

	// 1. 尝试从 MDD 获取音频
	if reader, contentType, err := s.getFromMDD(word); err == nil {
		return reader, contentType, nil
	}

	// 2. 尝试从独立 wav 文件获取
	if reader, contentType, err := s.getFromWavFile(word); err == nil {
		return reader, contentType, nil
	}

	// 3. 尝试从 LSA 获取
	if reader, contentType, err := s.getFromLSA(word); err == nil {
		return reader, contentType, nil
	}

	return nil, "", ErrAudioNotFound
}

// getFromMDD 从 MDD 字典资源中获取音频
func (s *AudioService) getFromMDD(word string) (io.Reader, string, error) {
	// 尝试常见的音频路径格式
	audioPatterns := []string{
		// 牛津词典格式（优先级最高，因为最常用）
		word + "__gb_1.mp3",
		word + "__gb_2.mp3",
		word + "__us_1.mp3",
		word + "__us_2.mp3",
		// 通用格式
		word + ".mp3",
		word + ".wav",
		word + ".ogg",
		"audio/" + word + ".mp3",
		"audio/" + word + ".wav",
		"sound/" + word + ".mp3",
		"sound/" + word + ".wav",
		strings.ToUpper(word[:1]) + "/" + word + ".mp3",
		strings.ToUpper(word[:1]) + "/" + word + ".wav",
	}

	for _, info := range s.mdxManager.ListLoaded() {
		if !info.HasMDD {
			continue
		}
		for _, pattern := range audioPatterns {
			if reader, err := s.mdxManager.GetResource(info.ID, pattern); err == nil {
				contentType := "audio/mpeg"
				if strings.HasSuffix(pattern, ".wav") {
					contentType = "audio/wav"
				} else if strings.HasSuffix(pattern, ".ogg") {
					contentType = "audio/ogg"
				}
				return reader, contentType, nil
			}
		}
	}

	return nil, "", ErrAudioNotFound
}

// getFromWavFile 从独立 wav 文件获取音频
// 文件路径格式：WyabdcRealPeopleTTS/{首字母}/{word}.wav
func (s *AudioService) getFromWavFile(word string) (io.Reader, string, error) {
	if len(word) == 0 {
		return nil, "", ErrAudioNotFound
	}

	// 构建文件路径
	firstLetter := strings.ToLower(string(word[0]))
	wavPath := filepath.Join(s.soundDir, "WyabdcRealPeopleTTS", firstLetter, word+".wav")

	data, err := os.ReadFile(wavPath)
	if err != nil {
		return nil, "", ErrAudioNotFound
	}

	return bytes.NewReader(data), "audio/wav", nil
}

// getFromLSA 从 LSA 音频包获取音频
func (s *AudioService) getFromLSA(word string) (io.Reader, string, error) {
	reader, err := s.ensureLSALoaded()
	if err != nil {
		return nil, "", ErrAudioNotFound
	}

	data, err := reader.Read(word)
	if err != nil {
		return nil, "", ErrAudioNotFound
	}

	return bytes.NewReader(data), "audio/wav", nil
}

// Close 关闭音频服务
func (s *AudioService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.lsaReader != nil {
		return s.lsaReader.Close()
	}
	return nil
}
