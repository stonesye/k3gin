package errors

import "fmt"

// ResponseError 定义响应错误信息
type ResponseError struct {
	Code    int
	Message string
	Status  int
	ERR     error
}

// Error 被实现  ResponError 也是 error
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

func WrapResponse(err error, code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args),
		Status:  status,
		ERR:     err,
	}

	return res
}

func Wrap400Response(err error, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    0,
		Message: fmt.Sprintf(msg, args),
		Status:  400,
		ERR:     err,
	}

	return res
}

func Wrap500Response(err error, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    0,
		Message: fmt.Sprintf(msg, args),
		Status:  400,
		ERR:     err,
	}

	return res
}

func NewResponse(code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args),
		Status:  status,
	}
	return res
}

func New400Response(msg string, args ...interface{}) error {
	return NewResponse(0, 400, msg, args...)
}

func New500Response(msg string, args ...interface{}) error {
	return NewResponse(0, 500, msg, args)
}
