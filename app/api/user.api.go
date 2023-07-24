package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/ginx"
	"k3gin/app/logger"
	"k3gin/app/schema"
	"k3gin/app/service"
	"net/http"
)

type UserApi struct {
	UserSrv *service.UserSrv
}

var UserApiSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

func (u *UserApi) Query(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = logger.NewTagContext(ctx, "__user__api__query__")
	var params schema.UserQueryParam

	// 将客户端提交的数据，封装到params中
	if err := ginx.ParseQuery(c, &params); err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}

	result, err := u.UserSrv.Query(ctx, params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Data})
}
