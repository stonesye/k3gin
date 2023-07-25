package errors

// 定义错误信息

var (
	ErrInternalServer = NewResponse(0, 500, "internal server error !")
	ErrNoPerm         = NewResponse(0, 401, "no permission")
	ErrNotFound       = NewResponse(0, 404, "not found")
	ErrTooManyRequest = NewResponse(0, 429, "too many request")
	ErrBadRequest     = New400Response("bad request")
	ErrInvalidParent  = New400Response("not found parent node ")
	ErrUserDisable    = New400Response("user forbidden")
)
