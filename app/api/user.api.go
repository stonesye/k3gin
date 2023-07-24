package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/ginx"
	"k3gin/app/schema"
	"k3gin/app/service"
)

type UserApi struct {
	UserSrv *service.UserSrv
}

var UserApiSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

// Query
// @Summary      Show user
// @Description  get userinfo by ID,UserName
// @Tags         QueryAPI
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  gin.H
// @Failure      400  {object}  schema.ErrorResult
// @Router       /api/v1/user [get]
func (u *UserApi) Query(c *gin.Context) {
	ctx := c.Request.Context()

	var params schema.UserQueryParam
	// 将客户端提交的数据，封装到params中
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := u.UserSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResList(c, result.Data)
}
