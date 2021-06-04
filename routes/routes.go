package routes

import (
	"github.com/gin-gonic/gin"
	"testNEW/settings"
)

func Init() {
	r := gin.Default()
	router := r.Group(settings.AppSettings.AppParams.ServerName)
	router.POST("/addserver", addServer)
	r.Run(":" + settings.AppSettings.AppParams.PortRun)
}
