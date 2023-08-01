package cron

import (
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"k3gin/app/cron/job"
)

type IWorker interface {
	Register(*v3cron.Cron)
}

var WorkerSet = wire.NewSet(wire.Struct(new(Worker), "*"), wire.Bind(new(IWorker), new(*Worker)))

type Worker struct {
	UserJob *job.UserJob
}

func (w *Worker) Register(cron *v3cron.Cron) {
	cron.AddJob(string(w.UserJob.Spec), v3cron.SkipIfStillRunning(v3cron.DefaultLogger)(w.UserJob)) // TODO 全部Cron都需要跳过当前正在执行的程序
}
