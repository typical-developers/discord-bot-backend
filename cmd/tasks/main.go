package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
	"github.com/typical-developers/discord-bot-backend/tasks"
)

type JobRegistryEntry struct {
	Name             string
	RunOnceAtStartup bool
	Disabled         bool
	Interval         string
	Task             func()
}

var JobRegistry = []JobRegistryEntry{
	{
		Name:             "Test",
		RunOnceAtStartup: true,
		Disabled:         true,
		Interval:         "*/1 * * * *",
		Task:             tasks.Test,
	},
	{
		Name:             "FlushWeeklyActivityLeaderboard",
		RunOnceAtStartup: false,
		Disabled:         false,
		Interval:         "0 0 * * 1",
		Task:             tasks.FlushWeeklyActivityLeaderboard,
	},
	{
		Name:             "FlushMonthlyActivityLeaderboard",
		RunOnceAtStartup: false,
		Disabled:         false,
		Interval:         "0 0 1 * *",
		Task:             tasks.FlushMonthlyActivityLeaderboard,
	},
}

var Cron = cron.New(cron.WithLocation(time.UTC))

func functionName(i interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	pieces := strings.Split(funcName, "/")

	return pieces[len(pieces)-1]
}

func main() {
	Cron.Start()

	for _, job := range JobRegistry {
		jobName := functionName(job.Task)

		if job.Disabled {
			logger.Log.Warn(fmt.Sprintf("%s disabled, skipping.", jobName))
			continue
		}

		_, err := Cron.AddFunc(job.Interval, job.Task)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("%s failed to register.", jobName), "error", err)
		}

		logger.Log.Info(fmt.Sprintf("%s successfully registered.", jobName))

		if job.RunOnceAtStartup {
			job.Task()
		}
	}

	runtime.Goexit()
}
