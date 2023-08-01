package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
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

type task struct {
	name string
	spec string
}

func (t *task) Run() {
	fmt.Println("name", t.name, "spec", t.spec)
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

	/**
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

	*/

	/**
	// 创建新号源， 控制cron的运行， 确保只有接触到特殊的信号以后， 主协程才会退出，子协程才会被回收
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 5; i++ {
		go func(data int) {
		EXIT:
			for {
				fmt.Println("开始协程", data)
				s := <-sig
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					fmt.Println("协程", data, "准备退出....")
					break EXIT
				case syscall.SIGHUP:
				default:
					fmt.Println("协程", data, "持续运行中....")
					time.Sleep(time.Second)
				}
			}
		}(i)
	}

	time.Sleep(10 * time.Second)

	*/

	// parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	// c := cron.New(cron.WithParser(parser))

	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(5)
	log.SetOutput(os.Stdout)
	log.WithField("aaaaa", "bbbbb")

	c := cron.New(cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log)))

	cron.WithChain()

	c.AddJob("* * * * * *", &task{
		name: "Test1",
		spec: "哈哈哈哈哈",
	})

	c.AddJob("* * * * * *", &task{
		name: "Test2",
		spec: "哈哈哈哈哈",
	})

	c.AddJob("@every 2s", &task{
		name: "Test3",
		spec: "哈哈哈哈哈",
	})

	c.Start()

	time.Sleep(20 * time.Second)

	c.Stop()
}
