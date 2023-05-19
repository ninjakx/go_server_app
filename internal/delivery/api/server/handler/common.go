package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// respondJSON makes the response with payload as json format
func respondJSON(c *gin.Context, status int, payload interface{}) (err error) {
	response, err := json.Marshal(payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.String(status, string(response))
	return
}

// respondError makes the error response with payload as json format
func respondError(c *gin.Context, code int, message string) {
	respondJSON(c, code, map[string]string{"error": message})
}
