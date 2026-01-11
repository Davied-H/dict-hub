package middleware

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SPAHandler 返回处理 SPA 路由的中间件
// 它会服务嵌入的静态文件，并将所有未匹配的路由回退到 index.html
func SPAHandler(staticFS fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(staticFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 路由和已知路径不处理，让其继续到下一个处理器
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/static/") ||
			strings.HasPrefix(path, "/dict-assets/") ||
			strings.HasPrefix(path, "/health") {
			c.Next()
			return
		}

		// 尝试提供静态文件
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		if _, err := fs.Stat(staticFS, filePath); err == nil {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// SPA 回退：返回 index.html
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
