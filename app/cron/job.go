package cron

// Job 包装一个对象，实现Run方法， 因为v3cron包里面的Job对象是需要实现Run函数的, 方便后续将任何待处理逻辑封装成Job，交给v3cron来处理
type Job struct {
	f func()
}

func (job *Job) Run() {
	job.f()
}

func NewJob(f func()) *Job {
	return &Job{f: f}
}
