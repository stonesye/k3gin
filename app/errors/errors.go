package errors

var (
	ErrInternalServer = NewResponse(0, 500, "internal server error !")
)
