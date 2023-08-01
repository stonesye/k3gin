package cron

import (
	"context"
)

// Context Cron执行的上下文，可以理解成为存储了一个任务的所有信息, 即单个任务
type Context struct {
	context.Context
	Name       string                 // 任务名
	Spec       string                 // 任务执行表达式
	Chain      []func(*Context)       // 任务要执行的函数
	ChainIndex int                    // 当前知道到第几个函数
	KV         map[string]interface{} // 存储KV
}

type IContext interface {
	Next()
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Value(interface{}) interface{}
}

func NewJobContext(name, spec string, funcs ...func(*Context)) *Context {
	return &Context{
		Context:    context.TODO(),
		Name:       name,
		Spec:       spec,
		Chain:      funcs,
		ChainIndex: 0,
		KV:         make(map[string]interface{}),
	}
}

// Next 执行Context中存储的下一个Chain
func (c *Context) Next() {
	curIndex := c.ChainIndex
	if curIndex >= len(c.Chain) { // 代表所有的chain都执行完了
		return
	}
	c.ChainIndex++
	// 执行chain中的函数
	c.Chain[curIndex](c)
}

func (c *Context) Set(s string, i interface{}) {
	c.KV[s] = i
}

func (c *Context) Get(s string) (interface{}, bool) {
	v, exist := c.KV[s]
	return v, exist
}

// Value job.Context 是包含了 context.Context 的, job.Context 拥有 context.Context 的函数实现的， 这里的Value只能说重写了
func (c *Context) Value(k interface{}) interface{} {
	if s, ok := k.(string); ok {
		v, _ := c.KV[s]
		return v
	}

	return c.Context.Value(k)
}

// FrameContext 框架Context 用于让job函数里面能获取到各种客户端来使用
type FrameContext struct {
	Ctx  *Context
	Cron *Cron
}
