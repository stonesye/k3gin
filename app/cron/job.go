package cron

import (
	"fmt"
	"k3gin/app/cron/context"
)

type Job struct {
	f func()
}

func newJob(f func()) *Job {
	return &Job{f: f}
}

func (j *Job) Run() {
	j.f()
}

func task1(ctx *context.CronContext) {

	fmt.Println("This is Task1")

}

func task2(ctx *context.CronContext) {
	fmt.Println("This is Task2")
}
