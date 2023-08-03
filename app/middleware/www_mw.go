package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

// StaticPathMiddleware 过滤请求过来的RequestURL 是进入静态目录地址还是进入应用程序内
// skippers 可以多个校验，目前这里只传了一个哟
func StaticPathMiddleware(root string, skippers ...func(*gin.Context) bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 判断请求的API是否合法, 如果是API，证明需要进入到应用逻辑所以直接到下一个Middleware 如果不是,就返回真实静态文件地址
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		fpath := filepath.Join(root, filepath.FromSlash(c.Request.URL.Path))

		// 判断文件是否存在，如果不存在就默认指向一个地址
		if _, err := os.Stat(fpath); err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, "index.html")
		}

		// 直接返回静态文件内容， 并终止程序继续往下执行, 因为静态地址是不需要进入其他逻辑的
		c.File(fpath)
		c.Abort()
	}
}
