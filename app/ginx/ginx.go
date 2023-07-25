package ginx

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"k3gin/app/errors"
	"k3gin/app/logger"
	"k3gin/app/schema"
	"net/http"
)

// 封装HTTP/HTTPS协议请求数据的封装和Response返回对象

const (
	REQBODYKEY = "req-body"
	RESBODYKEY = "res-body"
)

func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request query failed : %s", err.Error()))
	}

	return nil
}

func ParseJson(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request json failed : %s ", err.Error()))
	}

	return nil
}

func ResOK(c *gin.Context) {
	ResSuccess(c, schema.StatusResult{Status: "OK"})
}

func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.ListResult{List: v})
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

func ResJSON(ctx *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	ctx.Set(RESBODYKEY, buf)
	ctx.Data(status, "application/json; charset=utf-8", buf)
	ctx.Abort()
}

func ResError(ctx *gin.Context, err error, status ...int) {
	var res *errors.ResponseError

	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(status) > 0 {
		res.Status = status[0]
	}

	// 已封装好errors.ResponseErr, 但有可能error信息为空
	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}

		// 将err信息用日志记录下来
		if status := res.Status; status >= 400 && status < 500 {
			logger.WithContext(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.WithContext(logger.NewStackContext(ctx, err)).Errorf(err.Error())
		}
	}

	errItem := schema.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}

	ResJSON(ctx, res.Status, schema.ErrorResult{Error: errItem})

}
