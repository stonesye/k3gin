package main

import "fmt"

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
	c := New(WithName("lei"), WithFilePath("/config"), WithHost("127.0.0.1"))
	c.Run()
}
