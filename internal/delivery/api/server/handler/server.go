package handler

import (
	"GO_APP/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

// getServerOr404 gets a Server instance if exists, or respond the 404 error otherwise
func getServerOr404(db *gorm.DB, id int, c *gin.Context) (*model.Server, error) {
	server := model.Server{}
	err := db.Where("id = ?", id).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, err
}

func GetServerHostName(db *gorm.DB, c *gin.Context) {
	ps := c.Params
	thresh, err := strconv.Atoi(ps.ByName("thresh"))
	if err != nil {
		// pass default value
		thresh = DEFAULT_THESHOLD
		log.Printf("[server][GetServerHostName][strconv.Atoi] error:%+v\n", err)
		// respondError(w, http.StatusBadRequest, err.Error())
	}

	// Begin transaction
	tx := db.Begin()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	hostnames := []string{}
	err = tx.Table("servers").
		Select("hostname as Hostnames").
		Group("hostname").
		Having("COUNT(CASE WHEN active THEN 1 END) <= ?", thresh).
		Scan(&hostnames).Error
	tx.Commit()

	if err != nil {
		// tx.Rollback()
		log.Printf("[server][GetServerHostName][db.Table] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = respondJSON(c, http.StatusOK, hostnames)
	// Create log for the error
	if err != nil {
		log.Printf("[server][GetServerHostName][respondJSON] error:%+v\n", err)
	}
}

func CreateServer(db *gorm.DB, c *gin.Context) {
	server := model.Server{}
	r := c.Request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][CreateServer][decoder.Decode] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(&server).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[server][CreateServer][db.Create] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	err = respondJSON(c, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][CreateServer][respondJSON] error:%+v\n", err)
	}
}

func GetAllServer(db *gorm.DB, c *gin.Context) {
	servers := []model.Server{}
	err := db.Find(&servers).Error
	if err != nil {
		log.Printf("[server][GetAllServer][db.Find] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	err = respondJSON(c, http.StatusOK, servers)
	// Create log for the error
	if err != nil {
		log.Printf("[server][GetAllServer][respondJson] error:%+v\n", err)
	}
}

func GetServer(db *gorm.DB, c *gin.Context) {
	ps := c.Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][GetServer][strconv.Atoi] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	server, err := getServerOr404(db, id, c)
	if err != nil {
		log.Printf("[server][GetServer][getServerOr404] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	if server != nil {
		err = respondJSON(c, http.StatusOK, server)
		// Create log for the error
		if err != nil {
			log.Printf("[server][GetServer][respondJSON] error:%+v\n", err)
		}
	}

}

func UpdateServer(db *gorm.DB, c *gin.Context) {
	r := c.Request
	ps := c.Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][UpdateServer][strconv.Atoi] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	server, err := getServerOr404(db, id, c)
	if err != nil {
		log.Printf("[server][UpdateServer][getServerOr404] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&server); err != nil {
		log.Printf("[server][UpdateServer][decoder.Decode] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = tx.Model(&model.Server{}).Where("id = ?", server.ID).Updates(model.Server{IP: server.IP, Hostname: server.Hostname, Active: server.Active}).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[server][UpdateServer][tx.Commit] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	err = respondJSON(c, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][UpdateServer][respondJSON] error:%+v\n", err)
	}

}

func DisableServer(db *gorm.DB, c *gin.Context) {
	ps := c.Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][DisableServer][strconv.Atoi] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	server, err := getServerOr404(db, id, c)
	if err != nil {
		log.Printf("[server][DisableServer][getServerOr404] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	server.Disable()

	err = tx.Model(&model.Server{}).Where("id = ?", server.ID).Update("active", server.Active).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[server][DisableServer][tx.Commit] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	err = respondJSON(c, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][DisableServer][respondJSON] error:%+v\n", err)
	}

}

func EnableServer(db *gorm.DB, c *gin.Context) {
	ps := c.Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][EnableServer][strconv.Atoi] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	server, err := getServerOr404(db, id, c)
	if err != nil {
		log.Printf("[server][EnableServer][getServerOr404] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	server.Enable()

	err = tx.Model(&model.Server{}).Where("id = ?", server.ID).Update("active", server.Active).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[server][EnableServer][tx.Commit] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	err = respondJSON(c, http.StatusOK, server)
	// Create log for the error
	if err != nil {
		log.Printf("[server][EnableServer][respondJSON] error:%+v\n", err)
	}
}

func DeleteServer(db *gorm.DB, c *gin.Context) {
	ps := c.Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("[server][DeleteServer][strconv.Atoi] error:%+v\n", err)
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, err = getServerOr404(db, id, c)
	if err != nil {
		log.Printf("[server][DeleteServer][getServerOr404] error:%+v\n", err)
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	err = tx.Delete(&model.Server{}, id).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[server][DeleteServer][tx.Commit] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	err = respondJSON(c, http.StatusOK, nil)
	if err != nil {
		log.Printf("[server][DeleteServer][respondJSON] error:%+v\n", err)
		respondError(c, http.StatusInternalServerError, err.Error())
	}
}
