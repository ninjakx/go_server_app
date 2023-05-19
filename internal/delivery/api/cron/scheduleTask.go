package cron

import (
	"GO_APP/internal/delivery/api/cron/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SchedulerRoute struct {
	Router       *gin.Engine
	DB           *gorm.DB
	SchedulerJob *handler.Scheduler
}

// This will have server related api
// https://github.com/gin-gonic/gin/issues/1681
// Set all required routers

func (a *SchedulerRoute) SetSchedulerRouter() {
	router := a.Router
	// Routing for handling the projects
	router.POST("/scheduler/start", a.StartScheduler)
	router.POST("/scheduler/stop", a.StopScheduler)
}

// Handlers to start the scheduler
func (a *SchedulerRoute) StartScheduler(c *gin.Context) {
	a.SchedulerJob.StartSchedulerJob(c, a.DB)
}

// Handlers to stop the scheduler
func (a *SchedulerRoute) StopScheduler(c *gin.Context) {
	a.SchedulerJob.StopSchedulerJob(c)
}

// Run the SchedulerRoute on it's router
func (a *SchedulerRoute) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
