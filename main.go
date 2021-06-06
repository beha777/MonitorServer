package main

import (
	"MonitorServer/db"
	"MonitorServer/jobs"
	"MonitorServer/models"
	"MonitorServer/routes"
	"MonitorServer/server"
	"MonitorServer/settings"
	"time"
)

func main() {
	settings.AppSettings = settings.ReadSettings()
	db.ConnectDatabase()

	db.GetDBConn().DropTable(&models.Service{})
	db.GetDBConn().AutoMigrate(&models.Server{}, &models.ServerInfo{}, &models.Service{})

	var serversInfo []models.ServerInfo
	db.GetDBConn().Find(&serversInfo)
	for _, serverInfo := range serversInfo {
		db.AddServices(server.GetServicesList(server.ConnectToServer(serverInfo.ID)), serverInfo.ID)
	}
	go jobs.CheckServicesStart()
	go jobs.CheckServersStart()
	go jobs.CheckPingAndTelnet()
	routes.Init()
	time.Sleep(time.Minute)
}
