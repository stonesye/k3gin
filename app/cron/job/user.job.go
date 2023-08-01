package job

import (
	"fmt"
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
)

type UserJobName string
type UserJobSpec string

type UserJob struct {
	Name  UserJobName
	Spec  UserJobSpec
	Http  *httpx.Client
	DB    *gormx.DB
	Store *redisx.Store
}

var UserJobSet = wire.NewSet(wire.Struct(new(UserJob), "*"),
	wire.Value(UserJobName("user")),
	wire.Value(UserJobSpec("*/5 * * * * *")))

func (u *UserJob) Run() {
	fmt.Println(u.Name, "--------->正在运行")
}
