package middleware

import (
	"github.com/gin-gonic/gin"
)

// AllowPathPrefixSkipper 校验RequestURL 的连接是不是静态目录
func AllowPathPrefixSkipper(prefixes ...string) func(*gin.Context) bool {
	return func(c *gin.Context) bool {

		path := c.Request.URL.Path
		pathLen := len(path)

		for _, prefix := range prefixes {
			if pl := len(prefix); pathLen >= pl && path[:pl] == prefix {
				return true
			}
		}

		return false
	}
}

// AllowPathPrefixNoSkipper 凡是prefix的目录都返回true, 表示需要加traceID获取其他
func AllowPathPrefixNoSkipper(prefixes ...string) func(*gin.Context) bool {

	// prefixes 和 context.gin.Request.URL比对
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)
		for _, prefix := range prefixes {
			if pl := len(prefix); pathLen >= pl && path[:pl] == prefix {
				return false
			}
		}

		return true
	}
}

func SkipHandler(c *gin.Context, skippers ...func(ctx *gin.Context) bool) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}

	return false
}
