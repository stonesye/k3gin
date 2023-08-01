package job

import (
	"fmt"
)

type UserJob struct {
	Name string
}

func (u *UserJob) Run() {
	fmt.Println("do something for user : ", u.Name)
}
