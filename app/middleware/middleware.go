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

func SkipHandler(c *gin.Context, skippers ...func(ctx *gin.Context) bool) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}

	return false
}
