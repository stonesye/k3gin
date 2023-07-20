package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

// WWWMiddleware (root[静态目录地址] , AllowPathPrefixSkipper[允许访问的目录地址])
func WWWMiddleware(root string, skippers ...func(ctx *gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断API请求是否合法
		if SkipHandler(c, skippers...) {
			c.Next()
		}

		// 将请求的地址转移到静态文件中去
		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))

		if _, err := os.Stat(fpath); err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, "index.html")
		}

		c.File(fpath)
		c.Abort()
	}
}
