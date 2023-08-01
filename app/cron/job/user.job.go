package job

import (
	"fmt"
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"time"
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
	wire.Value(UserJobName("TestUserJob")),
	wire.Value(UserJobSpec("*/2 * * * * *")))

var count = 0

func (u *UserJob) Run() {
	count++
	fmt.Println(u.Name, "--------->正在运行", time.Now().Format("2006-01-02 15:04:05"))
	/**
	if count == 2 {
		panic("ooooooooooooooops !!!")
	}

	*/
	time.Sleep(5 * time.Second)

}
