package api

import (
	"GO_APP/config"
	"GO_APP/internal/handler"
	"GO_APP/internal/model"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// App has router and db instances
type App struct {
	Router *httprouter.Router //mux.Router
	DB     *sqlx.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBname,
	)

	db, err := sqlx.Connect(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	} else {
		fmt.Printf("Connected to database\n")
	}

	a.DB = model.DBMigrate(db)
	a.Router = httprouter.New()
	a.SetRouters()
}

// https://github.com/gin-gonic/gin/issues/1681
// Set all required routers
func (a *App) SetRouters() {
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
func (a *App) CreateServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.CreateServer(a.DB, w, r)
}

func (a *App) GetServerHostname(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.GetServerHostName(a.DB, w, ps)
}

func (a *App) GetServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.GetServer(a.DB, w, ps)
}

func (a *App) GetAllServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.GetAllServer(a.DB, w)
}

func (a *App) UpdateServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.UpdateServer(a.DB, w, r, ps)
}

func (a *App) DisableServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.DisableServer(a.DB, w, r, ps)
}

func (a *App) EnableServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.EnableServer(a.DB, w, r, ps)
}

func (a *App) DeleteServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handler.DeleteServer(a.DB, w, r, ps)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
