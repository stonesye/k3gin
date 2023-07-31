package cron

import (
	"fmt"
	"k3gin/app/config"
	"k3gin/app/cron/context"
	"reflect"
)

func task1(ctx context.CronContext) {

	fmt.Println("This is Task1")

}

func task2(ctx context.CronContext) {
	fmt.Println("This is Task2")
}

func ExecJob(ctx context.CronContext) {
	jobList := config.C.Cron.Jobs

	for _, job := range jobList {
		fmt.Println(job)
		funcVal := reflect.ValueOf(job.Task)
		funcVal.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}
