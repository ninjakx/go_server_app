package model

import (
	"GO_APP/internal/queries"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Id       int    `db:"id"` // not needed as auto increment just for making testing easier including it
	IP       string `db:"ip"`
	Hostname string `db:"hostname"`
	Active   bool   `db:"active"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *sqlx.DB) *sqlx.DB {
	db.MustExec(queries.CreateDB)
	return db
}

func (s *Server) Disable() {
	s.Active = false
}

func (s *Server) Enable() {
	s.Active = true
}
