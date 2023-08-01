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
	cron.AddJob(string(w.UserJob.Spec), w.UserJob)
}
