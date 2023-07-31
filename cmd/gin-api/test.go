package main

import (
	"context"
	"fmt"
	"time"
)

type Obj interface {
	Say()
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Say() {
	fmt.Println("Say ", p.Name)
}

type Man struct {
	Person
	Sex bool
}

func (m Man) Say() {
	fmt.Println("Say", m.Name)
	m.Person.Say()
}

func Test(o Obj) {
	var person Person
	var man Man
	if s, ok := o.(Person); ok {
		person = s
		person.Say()
	} else if s, ok := o.(Man); ok {
		man = s
		man.Say()
	}

}

func main() {
	fmt.Println("aaaaaa")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	fmt.Println("bbbbb")
	defer cancelFunc()
	fmt.Println("ctx", ctx)

	select {
	case <-ctx.Done():
		fmt.Println("收到信号")
	}
	fmt.Println(time.Now())
}
