package main

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	. "github.com/luckfire-go/cron-scheduler"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	_ "github.com/typical-developers/discord-bot-backend/internal/logger"
	"github.com/typical-developers/discord-bot-backend/services/cron/config"
	"github.com/typical-developers/discord-bot-backend/services/cron/tasks"
)

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d?%s",
		config.C.Database.Username,
		config.C.Database.Password,
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.Options,
	))

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	return db, nil
}

func main() {
	pqdb, err := dbConnect()
	if err != nil {
		panic(err)
	}
	queries := db.New(pqdb)
	tasks := tasks.NewTasks(pqdb, queries)

	registry := NewRegistry(cron.WithLocation(time.UTC))
	registry.OnJobAddSuccess = func(job *RegistryItem) {
		log.WithFields(log.Fields{
			"name": job.Name,
			"spec": job.Spec,
		}).Info("Successfully added job to registry.")
	}
	registry.OnJobAddFailure = func(job *RegistryItem, err error) {
		log.WithFields(log.Fields{
			"name": job.Name,
			"spec": job.Spec,
			"err":  err,
		}).Error("Failed to add job to registry.")
	}
	registry.OnJobFailed = func(job *RegistryItem, err error) {
		log.WithFields(log.Fields{
			"name": job.Name,
			"spec": job.Spec,
			"err":  err,
		}).Error("Job failed to run.")
	}

	registry.AddJobs([]RegistryItem{
		{
			Enabled:       true,
			RunOnRegister: true,

			Spec:     "0 0 * * 1",
			TaskFunc: tasks.FlushWeeklyActivityLeaderboard,
		},
		{
			Enabled:       true,
			RunOnRegister: true,

			Spec:     "0 0 1 * *",
			TaskFunc: tasks.FlushMonthlyActivityLeaderboard,
		},
	})

	registry.Start()
	runtime.Goexit()
}
