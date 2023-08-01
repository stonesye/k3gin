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

func Work(ctx context.Context, msg string) {
	for true {
		select {
		case <-ctx.Done():
			fmt.Println(msg, "goroutine is finish...")
			return
		default:
			fmt.Println("goroutine is running", time.Now().String())
			time.Sleep(time.Second)
		}
	}
}

func main() {
	/**
	ctx, cancel := context.WithCancel(context.Background())

	go Work(ctx, "withcancel")

	time.Sleep(time.Second * 3)

	fmt.Println("cancel ....")

	cancel()

	time.Sleep(time.Second * 3)

	fmt.Println("finish")

	*/

	/**
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
	*/

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 5; i++ {
		fmt.Println("start routine ", i)

		go func(ctx context.Context, data int) {

			go func(ctx context.Context, data int) {

				select {
				case <-ctx.Done():
					fmt.Println("协程", data, "的子协程关闭")

				}

			}(ctx, data)

			for true {
				select {
				case <-ctx.Done():
					fmt.Printf("协程 : %d 即将关闭\n", data)
					return
				default:
					fmt.Println(data, "协程在运行中....")
					time.Sleep(1 * time.Second)
				}

			}

		}(ctx, i)

	}

	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(10 * time.Second)

}
