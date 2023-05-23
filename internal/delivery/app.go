package api

import (
	"GO_APP/config"
	"GO_APP/internal/delivery/api/cron"
	"GO_APP/internal/delivery/api/cron/handler"
	"GO_APP/internal/delivery/api/server"
	"GO_APP/internal/delivery/api/user"
	"GO_APP/internal/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	ServiceRouter   server.ServerRoute
	DB              *gorm.DB
	SchedulerRouter cron.SchedulerRoute
	UserAuthRouter  user.UserAuthRoute
}

// App initialize with predefined configuration
func (a *App) Init(config *config.Config) {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBname,
	)

	// db, err := sqlx.Connect(config.DB.Dialect, dbURI)
	db, err := gorm.Open(postgres.Open(dbURI))
	if err != nil {
		log.Fatal("Could not connect database")
	} else {
		log.Printf("Connected to database\n")
	}

	a.DB = model.DBMigrate(db)

	// set service routers
	// serviceRouter := a.ServiceRouter
	eng := gin.New()
	a.ServiceRouter.Router = eng
	a.ServiceRouter.DB = a.DB
	a.ServiceRouter.SetServiceRouter()

	a.SchedulerRouter.Router = gin.New()
	a.SchedulerRouter.DB = a.DB
	a.SchedulerRouter.SchedulerJob = handler.InitializeScheduler()
	a.SchedulerRouter.SetSchedulerRouter()

	a.UserAuthRouter.Router = eng
	a.UserAuthRouter.DB = a.DB
	a.UserAuthRouter.SetUserAuthRoute()

}

// Run the app on it's router
func (a *App) Run(host string) {
	a.ServiceRouter.Run(host)
}

func (a *App) RunCron(host string) {
	a.SchedulerRouter.Run(host)
}
