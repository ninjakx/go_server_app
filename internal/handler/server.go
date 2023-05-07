package handler

import (
	"GO_APP/internal/model"
	"GO_APP/internal/queries"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// getServerOr404 gets a Server instance if exists, or respond the 404 error otherwise
func getServerOr404(db *sqlx.DB, id int, w http.ResponseWriter) (*model.Server, error) {
	server := model.Server{}
	err := db.Get(&server, queries.QueryFindServer, id)
	if err != nil {
		log.Printf("[server][getServerOr404][db.Get] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}
	return &server, err
}

func GetServerHostName(db *sqlx.DB, w http.ResponseWriter, ps httprouter.Params) {
	hostnames := []string{}
	thresh, err := strconv.Atoi(ps.ByName("thresh"))
	if err != nil {
		// pass default value
		thresh = DEFAULT_THESHOLD
		log.Printf("[server][GetServerHostName][strconv.Atoi] error:%+v\n", err)
		// respondError(w, http.StatusBadRequest, err.Error())
	}
	tx := db.MustBegin()
	tx.Select(&hostnames, queries.QueryGetAllHostnameWithThresh, thresh)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][GetServerHostName][tx.Commit] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	err = respondJSON(w, http.StatusOK, hostnames)
	// Create log for the error
	if err != nil {
		log.Printf("[server][GetServerHostName][respondJSON] error:%+v\n", err)
		return
	}

}

func CreateServer(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	server := model.Server{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][CreateServer][decoder.Decode] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}

	// Begin transaction
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("[server][CreateServer][db.Beginx] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	tx.NamedExec(queries.QueryInsertServerData, &server)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][CreateServer][tx.Commit()] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	err = respondJSON(w, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][CreateServer][respondJSON] error:%+v\n", err)
		return
	}
}

func GetAllServer(db *sqlx.DB, w http.ResponseWriter) {
	servers := []model.Server{}
	err := db.Select(&servers, queries.QueryAllserver)
	if err != nil {
		log.Printf("[server][GetAllServer][db.Select] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	err = respondJSON(w, http.StatusOK, servers)
	// Create log for the error
	if err != nil {
		log.Printf("[server][GetAllServer][respondJson] error:%+v\n", err)
		return
	}

}

func GetServer(db *sqlx.DB, w http.ResponseWriter, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][GetServer][strconv.Atoi] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}

	server, err := getServerOr404(db, id, w)
	if err != nil {
		log.Printf("[server][GetServer][getServerOr404] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	err = respondJSON(w, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][GetServer][respondJSON] error:%+v\n", err)
		return
	}

}

func UpdateServer(db *sqlx.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][UpdateServer][strconv.Atoi] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}

	// Begin transaction
	tx, err := db.Beginx()
	defer tx.Rollback()
	if err != nil {
		log.Printf("[server][UpdateServer][db.Beginx] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	server, err := getServerOr404(db, id, w)
	if err != nil {
		log.Printf("[server][UpdateServer][getServerOr404] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][UpdateServer][decoder.Decode] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	tx.NamedExec(queries.QueryUpdateServer, &server)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][UpdateServer][tx.Commit] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	err = respondJSON(w, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][UpdateServer][respondJSON] error:%+v\n", err)
		return
	}

}

func DisableServer(db *sqlx.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][DisableServer][strconv.Atoi] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	// Begin transaction
	tx, err := db.Beginx()
	defer tx.Rollback()
	if err != nil {
		log.Printf("[server][DisableServer][db.Beginx] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	server, err := getServerOr404(db, id, w)
	if err != nil {
		log.Printf("[server][DisableServer][getServerOr404] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	server.Disable()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][DisableServer][decoder.Decode] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	tx.NamedExec(queries.QueryUpdateServerStatus, &server)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][DisableServer][tx.Commit] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	err = respondJSON(w, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][DisableServer][respondJSON] error:%+v\n", err)
		return
	}

}

func EnableServer(db *sqlx.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][EnableServer][strconv.Atoi] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	// Begin transaction
	tx, err := db.Beginx()
	defer tx.Rollback()
	if err != nil {
		log.Printf("[server][EnableServer][db.Beginx] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	server, err := getServerOr404(db, id, w)
	if err != nil {
		log.Printf("[server][EnableServer][getServerOr404] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	server.Enable()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][EnableServer][decoder.Decode] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	tx.NamedExec(queries.QueryUpdateServerStatus, &server)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][EnableServer][tx.Commit] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	err = respondJSON(w, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][EnableServer][respondJSON] error:%+v\n", err)
		return
	}

}

func DeleteServer(db *sqlx.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][DeleteServer][strconv.Atoi] error:%+v\n", err)
		respondError(w, http.StatusBadRequest, err.Error())
	}
	// Begin transaction
	tx, err := db.Beginx()
	defer tx.Rollback()
	if err != nil {
		log.Printf("[server][DeleteServer][db.Beginx] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	_, err = getServerOr404(db, id, w)
	if err != nil {
		log.Printf("[server][DeleteServer][getServerOr404] error:%+v\n", err)
		respondError(w, http.StatusNotFound, err.Error())
	}

	tx.MustExec(queries.QueryDeleteServer, id)
	err = tx.Commit()
	if err != nil {
		log.Printf("[server][DeleteServer][tx.Commit] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	err = respondJSON(w, http.StatusOK, nil)
	if err != nil {
		log.Printf("[server][DeleteServer][respondJSON] error:%+v\n", err)
		respondError(w, http.StatusInternalServerError, err.Error())
	}
}
