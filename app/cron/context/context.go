package context

import "context"

type HandleFunc func(*CronContext)

type CronContext struct { // 每个Job的独立Context
	context.Context                        // 系统的Context, 匿名到自定义的Context, 继承了系统Context的所有函数
	Name            string                 // 任务名
	Spec            string                 // 任务表达式
	Chain           []HandleFunc           // 中间件调用链
	ChainIndex      int                    // 当前chain执行到哪里
	KeyValues       map[string]interface{} // 存储KV
}

// NewCronContext 给每个JOB生成一个新的context
func NewCronContext(name string, spec string, middlewares []HandleFunc) *CronContext {
	return &CronContext{
		Context:    context.TODO(),
		Name:       name,
		Spec:       spec,
		Chain:      middlewares,
		ChainIndex: 0,
		KeyValues:  make(map[string]interface{}, 0),
	}
}

// Next 模拟调用下一个中间件middleware
func (c *CronContext) Next() {
	curIndex := c.ChainIndex
	if curIndex >= len(c.Chain) { // 没有下一个middleware
		return
	}

	c.ChainIndex++
	c.Chain[curIndex](c) // 执行下一个middleware
}

func (c *CronContext) Set(k string, v interface{}) {
	c.KeyValues[k] = v
}

func (c *CronContext) Get(k string) (v interface{}, exist bool) {
	v, exist = c.KeyValues[k]
	return
}

func (c *CronContext) Value(k interface{}) (v interface{}) { // 重写context.Context 的value
	if s, ok := k.(string); ok {
		v, _ = c.KeyValues[s]
		return v
	}
	return c.Context.Value(k)
}
