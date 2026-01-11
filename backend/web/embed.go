package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var staticFiles embed.FS

// GetFS 返回嵌入的前端静态文件系统
// 在 Docker 构建时，前端构建产物会被复制到 web/dist 目录
// 然后通过 go:embed 嵌入到二进制文件中
func GetFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "dist")
}
