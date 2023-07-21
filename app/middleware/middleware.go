package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AllowPathPrefixSkipper(prefixes ...string) func(*gin.Context) bool {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		fmt.Println(path, "<---path--->")

		pathLen := len(path)

		for _, prefix := range prefixes {
			fmt.Println(prefix, "<---prefix--->")
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
