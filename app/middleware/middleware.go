package middleware

import "github.com/gin-gonic/gin"

// AllowPathPrefixSkipper 判断是不是可以放过的API请求URL
func AllowPathPrefixSkipper(prefixes ...string) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}

		return false
	}
}

// AllowPathPrefixNoSkipper 判断哪些是不可以进入系统的API请求URL
func AllowPathPrefixNoSkipper(prefixes ...string) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}

		return true
	}

}

func SkipHandler(c *gin.Context, skippers ...func(*gin.Context) bool) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}
