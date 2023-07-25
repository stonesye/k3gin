package schema

type OrderDirection int

type OrderField struct { // 排序字段
	Key       string         // 排序字段
	Direction OrderDirection // 升序还是降序
}

// ErrorResult 错误返回给客户端
type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

// SuccessResult 成功返回给客户端
type SuccessResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
