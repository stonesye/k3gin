package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

// RouterSet 利用wire设置Router构造函数, 细心的发型当前文件并没有函数是可以构造Router Struct的， 所以利用wire.Struct来构造，这里特殊的事Router实现了接口，我们还需要利用wire.Bind绑定下接口
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(engine *gin.Engine) error
	Prefixes() []string
}

type Router struct {
}

// Prefixes API允许访问的目录地址
func (r *Router) Prefixes() []string {
	return []string{"/api"}
}

func (r *Router) Register(e *gin.Engine) error {
	r.RegisterAPI(e)
	return nil
}

func (r *Router) RegisterAPI(e *gin.Engine) {

	g := e.Group("/api")

	v1 := g.Group("/v1")
	{
		guser := v1.Group("/user")
		{
			guser.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "pong"})
			})
		}
	}

}
