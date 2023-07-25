package errors

import "fmt"

// 封装 Response 的数据

// ResponseError 定义响应信息
type ResponseError struct {
	Code    int
	Message string
	Status  int
	ERR     error
}

// Error接口被实现  *ResponseError 也是 error
func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

// UnWrapResponse 判断error是不是ResponseError
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}

	return nil
}

// WrapResponse 封装ResponseError
func WrapResponse(err error, code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args),
		Status:  status,
		ERR:     err,
	}

	return res
}

// Wrap400Response 封装400状态码的ResponseError
func Wrap400Response(err error, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    0,
		Message: fmt.Sprintf(msg, args),
		Status:  400,
		ERR:     err,
	}

	return res
}

// Wrap500Response 封装500状态码的ResponseError
func Wrap500Response(err error, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    0,
		Message: fmt.Sprintf(msg, args),
		Status:  400,
		ERR:     err,
	}

	return res
}

// NewResponse 封装一个任意状态码的ResponseError, 但ResponseError.ERR = nil
func NewResponse(code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args),
		Status:  status,
	}
	return res
}

// New400Response 封装400状态码的ResponseError, 但ResponseError.ERR = nil
func New400Response(msg string, args ...interface{}) error {
	return NewResponse(0, 400, msg, args...)
}

// New500Response 封装500状态码的ResponseError, 但ResponseError.ERR = nil
func New500Response(msg string, args ...interface{}) error {
	return NewResponse(0, 500, msg, args)
}
