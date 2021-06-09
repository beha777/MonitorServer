package routes

import (
	"MonitorServer/settings"
	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	router := r.Group(settings.AppSettings.AppParams.ServerName)
	router.POST("/addserver", addServer)
	r.Run(":" + settings.AppSettings.AppParams.PortRun)
}
