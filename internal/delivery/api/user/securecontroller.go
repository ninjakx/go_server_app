package user

import (
	middlewares "GO_APP/internal/delivery/api/middleware"
	"GO_APP/internal/delivery/api/user/controller"
	"GO_APP/internal/delivery/api/user/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAuthRoute struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (a *UserAuthRoute) SetUserAuthRoute() {
	router := a.Router
	api := router.Group("/user/auth")
	{
		// Routing for handling the projects
		api.POST("/token", a.GenerateToken)
		api.POST("/user/register", a.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", handler.Ping)
		}
	}
}

func (a *UserAuthRoute) GenerateToken(c *gin.Context) {
	controller.GenerateToken(a.DB, c)
}
func (a *UserAuthRoute) RegisterUser(c *gin.Context) {
	handler.RegisterUser(a.DB, c)
}

// Run the UserAuthRoute on it's router
func (a *UserAuthRoute) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
