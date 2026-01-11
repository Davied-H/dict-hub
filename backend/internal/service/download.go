package service

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"dict-hub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrInvalidURL        = errors.New("invalid URL")
	ErrDownloadFailed    = errors.New("download failed")
	ErrTaskNotFound      = errors.New("download task not found")
	ErrInvalidFileFormat = errors.New("invalid file format, only .mdx files are supported")
)

// DownloadService 下载服务
type DownloadService struct {
	db            *gorm.DB
	dictDir       string
	dictSourceSvc *DictSourceService
	client        *http.Client
	mu            sync.Mutex
}

// NewDownloadService 创建下载服务
func NewDownloadService(db *gorm.DB, dictDir string, dictSourceSvc *DictSourceService) *DownloadService {
	return &DownloadService{
		db:            db,
		dictDir:       dictDir,
		dictSourceSvc: dictSourceSvc,
		client:        &http.Client{},
	}
}

// StartDownload 启动异步下载任务
func (s *DownloadService) StartDownload(downloadURL string) (*model.DownloadTask, error) {
	// 验证 URL
	parsedURL, err := url.Parse(downloadURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, ErrInvalidURL
	}

	// 从 URL 提取文件名
	fileName := filepath.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		fileName = "dictionary.mdx"
	}

	// 确保是 .mdx 文件
	if !strings.HasSuffix(strings.ToLower(fileName), ".mdx") {
		return nil, ErrInvalidFileFormat
	}

	// 创建下载任务
	task := &model.DownloadTask{
		URL:      downloadURL,
		FileName: fileName,
		Status:   model.DownloadStatusPending,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// 启动后台下载
	go s.downloadWorker(task.ID)

	return task, nil
}

// GetTaskStatus 获取下载任务状态
func (s *DownloadService) GetTaskStatus(taskID uint) (*model.DownloadTask, error) {
	var task model.DownloadTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

// downloadWorker 后台下载工作协程
func (s *DownloadService) downloadWorker(taskID uint) {
	var task model.DownloadTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return
	}

	// 更新状态为下载中
	s.updateTaskStatus(taskID, model.DownloadStatusDownloading, 0, 0, 0, "")

	// 发起 HTTP 请求
	resp, err := s.client.Get(task.URL)
	if err != nil {
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "failed to connect: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "HTTP error: "+resp.Status)
		return
	}

	totalSize := resp.ContentLength

	// 确保目录存在
	if err := os.MkdirAll(s.dictDir, 0755); err != nil {
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "failed to create directory: "+err.Error())
		return
	}

	// 生成唯一文件名（避免冲突）
	filePath := filepath.Join(s.dictDir, task.FileName)
	filePath = s.uniqueFilePath(filePath)

	// 创建临时文件
	tmpPath := filePath + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "failed to create file: "+err.Error())
		return
	}

	// 下载并更新进度
	var downloaded int64
	buf := make([]byte, 32*1024)
	lastProgress := 0

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				file.Close()
				os.Remove(tmpPath)
				s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "failed to write file: "+writeErr.Error())
				return
			}
			downloaded += int64(n)

			// 计算进度
			progress := 0
			if totalSize > 0 {
				progress = int(float64(downloaded) / float64(totalSize) * 100)
			}

			// 每5%更新一次进度
			if progress >= lastProgress+5 || progress == 100 {
				s.updateTaskStatus(taskID, model.DownloadStatusDownloading, progress, totalSize, downloaded, "")
				lastProgress = progress
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			os.Remove(tmpPath)
			s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "download interrupted: "+err.Error())
			return
		}
	}
	file.Close()

	// 重命名临时文件
	if err := os.Rename(tmpPath, filePath); err != nil {
		os.Remove(tmpPath)
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 0, 0, 0, "failed to save file: "+err.Error())
		return
	}

	// 添加到字典源
	dictSource, err := s.dictSourceSvc.Add(filePath)
	if err != nil {
		s.updateTaskStatus(taskID, model.DownloadStatusFailed, 100, totalSize, downloaded, "file saved but failed to load dictionary: "+err.Error())
		return
	}

	// 更新任务为完成
	s.db.Model(&model.DownloadTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":         model.DownloadStatusCompleted,
		"progress":       100,
		"total_size":     totalSize,
		"download_size":  downloaded,
		"dict_source_id": dictSource.ID,
	})
}

// updateTaskStatus 更新任务状态
func (s *DownloadService) updateTaskStatus(taskID uint, status string, progress int, totalSize, downloadSize int64, errorMsg string) {
	updates := map[string]interface{}{
		"status":   status,
		"progress": progress,
	}
	if totalSize > 0 {
		updates["total_size"] = totalSize
	}
	if downloadSize > 0 {
		updates["download_size"] = downloadSize
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	s.db.Model(&model.DownloadTask{}).Where("id = ?", taskID).Updates(updates)
}

// uniqueFilePath 生成唯一文件路径
func (s *DownloadService) uniqueFilePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	for i := 1; i < 1000; i++ {
		newPath := base + "_" + string(rune('0'+i/100)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10)) + ext
		if i < 10 {
			newPath = base + "_" + string(rune('0'+i)) + ext
		} else if i < 100 {
			newPath = base + "_" + string(rune('0'+i/10)) + string(rune('0'+i%10)) + ext
		}
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}

	return path
}
