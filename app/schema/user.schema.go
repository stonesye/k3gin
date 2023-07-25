package schema

import (
	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       uint64 `json:"id,string"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Status   int    `json:"status"`
}

func (u *User) String() string {
	s, err := jsoniter.MarshalToString(u)
	if err != nil {
		return ""
	}
	return s
}

func (u *User) CleanPassword() *User {
	u.Password = ""
	return u
}

type UserQueryParam struct { // 查询条件
	UserName   string `form:"user_name"`   // 用户名
	QueryValue string `form:"query_value"` // 模糊查询
	Status     int    `form:"status"`      // 用户状态
}

type UserQueryOptions struct {
	OrderFields  []*OrderField // 需要order的字段
	SelectFields []string      // 需要select的字段
}

type Users []*User

type UserQueryResult struct { // 查询结果
	Data Users
}
