package job

import "fmt"

type User struct {
	Name string
	Spec string
}

func (u *User) Run() {
	fmt.Println("user do task .....")
}
