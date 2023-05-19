package api

import (
	"GO_APP/config"
	"GO_APP/internal/delivery/api/cron"
	"GO_APP/internal/delivery/api/cron/handler"
	"GO_APP/internal/delivery/api/server"
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
		fmt.Printf("Connected to database\n")
	}

	a.DB = model.DBMigrate(db)

	// set service routers
	// serviceRouter := a.ServiceRouter
	a.ServiceRouter.Router = gin.New()
	a.ServiceRouter.DB = a.DB
	a.ServiceRouter.SetServiceRouter()

	a.SchedulerRouter.Router = gin.New()
	a.SchedulerRouter.DB = a.DB
	a.SchedulerRouter.SchedulerJob = handler.InitializeScheduler()
	a.SchedulerRouter.SetSchedulerRouter()
}

// // https://github.com/gin-gonic/gin/issues/1681
// // Set all required routers
// func (a *App) SetRouters() {
// 	router := a.Router
// 	// Routing for handling the projects
// 	router.GET("/servers/get_hostname/:thresh", a.GetServerHostname)
// 	router.GET("/servers", a.GetAllServer)
// 	router.GET("/server/:id", a.GetServer)
// 	router.POST("/servers/create", a.CreateServer)
// 	router.PUT("/servers/:id/update_server", a.UpdateServer)
// 	router.PUT("/servers/:id/disable", a.DisableServer)
// 	router.PUT("/servers/:id/enable", a.EnableServer)
// 	router.DELETE("/servers/:id", a.DeleteServer)
// }

// // Handlers to manage Server Data
// func (a *App) CreateServer(c *gin.Context) {
// 	handler.CreateServer(a.DB, c)
// }

// func (a *App) GetServerHostname(c *gin.Context) {
// 	handler.GetServerHostName(a.DB, c)
// }

// func (a *App) GetServer(c *gin.Context) {
// 	handler.GetServer(a.DB, c)
// }

// func (a *App) GetAllServer(c *gin.Context) {
// 	handler.GetAllServer(a.DB, c)
// }

// func (a *App) UpdateServer(c *gin.Context) {
// 	handler.UpdateServer(a.DB, c)
// }

// func (a *App) DisableServer(c *gin.Context) {
// 	handler.DisableServer(a.DB, c)
// }

// func (a *App) EnableServer(c *gin.Context) {
// 	handler.EnableServer(a.DB, c)
// }

// func (a *App) DeleteServer(c *gin.Context) {
// 	handler.DeleteServer(a.DB, c)
// }

// Run the app on it's router
func (a *App) Run(host string) {
	a.ServiceRouter.Run(host)
	// log.Fatal(http.ListenAndServe(host, a.ServiceRouter.Router))
}

func (a *App) RunCron(host string) {
	a.SchedulerRouter.Run(host)
}
