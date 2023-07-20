package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Name     string
	FilePath string
	Host     string
}

func WithHost(s string) func(*Config) {
	return func(c *Config) {
		c.Host = s
	}
}

func WithName(name string) func(*Config) {
	return func(c *Config) {
		c.Name = name
	}
}

func WithFilePath(s string) func(*Config) {
	return func(c *Config) {
		c.FilePath = s
	}
}

func New(opt ...func(*Config)) *Config {
	var c Config
	for _, o := range opt {
		o(&c)
	}
	return &c
}

func (c *Config) Run() {
	fmt.Println(c.FilePath, c.Name, c.Host)
}

func main() {
	/**
	c := New(WithName("lei"), WithFilePath("/config"), WithHost("127.0.0.1"))
	c.Run()
	*/

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, request *http.Request) {
		fmt.Println("got / request")
		io.WriteString(resp, "this is my web\n")
	})

	mux.HandleFunc("/go", func(resp http.ResponseWriter, request *http.Request) {
		fmt.Println("go /go request")
		io.WriteString(resp, "this is go reqeust")
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:8081",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		fmt.Println("开始监听.....")
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}

		fmt.Println("结束监听....")

	}()

	fmt.Println("主协程......1")
	time.Sleep(time.Second * 20)
	fmt.Println("住协程......2")

}
