package routes

import (
	"MonitorServer/db"
	"MonitorServer/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addServer(context *gin.Context) {
	var server models.ServerInfo
	err := json.NewDecoder(context.Request.Body).Decode(server)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"response": err.Error(),
		})
		return
	}
	if db.AddServer(server) {
		context.JSON(http.StatusOK, gin.H{
			"response": "Server successfully added",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"response": "Server exists",
		})
	}
}
