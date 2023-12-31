package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/ginx"
	"k3gin/app/httpx"
	"k3gin/app/schema"
	"k3gin/app/service"
	"time"
)

type UserApi struct {
	UserSrv *service.UserSrv
	Client  *httpx.Client
}

var UserApiSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

// Query
//
//	@Summary	根据用户名或用户状态查询用户信息
//	@Tags		UserQueryAPI
//	@Param		user_name	query		string				false	"用户名"
//	@Param		status		query		int					false	"用户状态(1，正常; 2，失效)"
//	@Param		query_value query		string				false	"模糊查询"
//	@Success	200			{object}	schema.SuccessResult "用户列表"
//	@Failure	400			{object}	schema.ErrorResult	"错误信息"
//	@Router		/api/v1/user [get]
func (u *UserApi) Query(c *gin.Context) {
	ctx := c.Request.Context()

	var params schema.UserQueryParam
	// 将客户端提交的数据，封装到params中
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	// 排序和查询字段
	options := schema.UserQueryOptions{
		OrderFields: []*schema.OrderField{
			{Key: "user_id", Direction: 1},
		},
		SelectFields: []string{"user_id", "username", "password", "status"},
	}

	result, err := u.UserSrv.Query(ctx, params, options)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result.Data, "[%v]user list success !", time.Now())
}

func (u *UserApi) Get(c *gin.Context) {
	/**
	go func() {
		res := u.Client.Get("http://127.0.0.1:8081/api/v1/user", nil, nil)
		time.Sleep(3 * time.Second)
		logger.WithContext(c.Request.Context()).Infof("请求信息:%v", res)
	}()

	*/
	ginx.ResSuccess(c, "success", "ok")
}
