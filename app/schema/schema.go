package schema

type StatusResult struct {
	Status string `json:"status"`
}

type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResult struct {
	Error ErrorItem `json:"error"`
}

type ListResult struct {
	List interface{} `json:"list"`
}

type OrderDirection int

type OrderField struct { // 排序字段
	Key       string         // 排序字段
	Direction OrderDirection // 升序还是降序
}
