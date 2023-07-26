package ginx

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"k3gin/app/contextx"
	"k3gin/app/errors"
	"k3gin/app/logger"
	"k3gin/app/schema"
	"net/http"
)

// 封装HTTP/HTTPS协议请求数据的封装和Response返回对象

const (
	ReqBodyKey  = "req-body"
	ResBodyKey  = "res-body"
	SuccessCode = 0
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

func ResSuccess(c *gin.Context, v interface{}, format string, option ...interface{}) {
	res := schema.SuccessResult{
		Code:    SuccessCode,
		Message: fmt.Sprintf(format, option),
		Data:    v,
	}
	ResJson(c, http.StatusOK, res)
}

func ResError(c *gin.Context, err error, httpStatus ...int) {
	var res *errors.ResponseError
	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else { // 如果不是ResponseError
			res = errors.UnWrapResponse(errors.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(httpStatus) > 0 {
		res.Status = httpStatus[0]
	}

	// 记录错误日志
	if err := res.ERR; err != nil {
		// 将err信息用日志记录下来
		if status := res.Status; status >= 400 && status < 500 {
			logger.WithContext(c).Warnf(err.Error())
		} else if status >= 500 {
			logger.WithContext(contextx.NewStack(c, err)).Errorf(err.Error())
		}
	}

	ResJson(c, res.Status, schema.ErrorResult{
		Code:    res.Code,
		Message: res.Message,
		Err:     res.ERR,
	})
}

func ResJson(c *gin.Context, httpStatus int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(httpStatus, "application/json; charset=utf-8", buf)
	c.Abort()
}
