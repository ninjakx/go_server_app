package server

import (
	"GO_APP/internal/delivery/api/server/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServerRoute struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// This will have server related api
// https://github.com/gin-gonic/gin/issues/1681
// Set all required routers

func (a *ServerRoute) SetServiceRouter() {
	router := a.Router
	// Routing for handling the projects
	router.GET("/servers/get_hostname/:thresh", a.GetServerHostname)
	router.GET("/servers", a.GetAllServer)
	router.GET("/server/:id", a.GetServer)
	router.POST("/servers/create", a.CreateServer)
	router.PUT("/servers/:id/update_server", a.UpdateServer)
	router.PUT("/servers/:id/disable", a.DisableServer)
	router.PUT("/servers/:id/enable", a.EnableServer)
	router.DELETE("/servers/:id", a.DeleteServer)
}

// Handlers to manage Server Data
func (a *ServerRoute) CreateServer(c *gin.Context) {
	handler.CreateServer(a.DB, c)
}

func (a *ServerRoute) GetServerHostname(c *gin.Context) {
	handler.GetServerHostName(a.DB, c)
}

func (a *ServerRoute) GetServer(c *gin.Context) {
	handler.GetServer(a.DB, c)
}

func (a *ServerRoute) GetAllServer(c *gin.Context) {
	handler.GetAllServer(a.DB, c)
}

func (a *ServerRoute) UpdateServer(c *gin.Context) {
	handler.UpdateServer(a.DB, c)
}

func (a *ServerRoute) DisableServer(c *gin.Context) {
	handler.DisableServer(a.DB, c)
}

func (a *ServerRoute) EnableServer(c *gin.Context) {
	handler.EnableServer(a.DB, c)
}

func (a *ServerRoute) DeleteServer(c *gin.Context) {
	handler.DeleteServer(a.DB, c)
}

// Run the ServerRoute on it's router
func (a *ServerRoute) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
