package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	job       *gocron.Job
}

func (sch *Scheduler) StartSchedulerJob(c *gin.Context, db *gorm.DB) {
	if sch == nil {
		log.Println("Scheduler not initialized")
		return
	}

	// If scheduler is not null and scheduler is running then create a new job
	if sch.job == nil || !sch.scheduler.IsRunning() {
		sch.job, _ = sch.scheduler.Every(2).Second().Do(func() {
			get_hostname(db)
		})
		sch.scheduler.StartAsync()

		c.String(http.StatusOK, "Cron job started")
	} else {
		c.String(http.StatusOK, "Cron job is already running")
	}

}

func (sch *Scheduler) StopSchedulerJob(c *gin.Context) {
	if sch.scheduler != nil && sch.scheduler.IsRunning() {
		sch.scheduler.RemoveByReference(sch.job)
		sch.job = nil
		c.String(http.StatusOK, "Cron job stopped")
	} else {
		c.String(http.StatusOK, "No active cron job to stop")
	}
}

func InitializeScheduler() *Scheduler {
	sch := gocron.NewScheduler(time.Local)
	return &Scheduler{
		scheduler: sch,
	}
}
