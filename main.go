package main

import (
	"MonitorServer/db"
	"MonitorServer/jobs"
	"MonitorServer/models"
	"MonitorServer/routes"
	"MonitorServer/settings"
	"time"
)

func main() {
	settings.AppSettings = settings.ReadSettings()

	db.ConnectDatabase()

	//db.GetDBConn().DropTable(&models.Service{})
	db.GetDBConn().AutoMigrate(&models.Server{}, &models.ServerInfo{}, &models.Service{})
	db.GetDBConn().Updates(&models.Server{
		LastTime:     settings.AppSettings.PeriodParams.NilTime,
		LastNotified: settings.AppSettings.PeriodParams.NilTime,
	})
	db.GetDBConn().Updates(&models.Service{
		LastTime:     settings.AppSettings.PeriodParams.NilTime,
		LastNotified: settings.AppSettings.PeriodParams.NilTime,
	})
	var serversInfo []models.ServerInfo
	db.GetDBConn().Find(&serversInfo)
	/*for _, serverInfo := range serversInfo {
		db.AddServices(server.GetServicesList(server.ConnectToServer(serverInfo.ID)), serverInfo.ID)
	}*/
	go jobs.CheckServicesStart()
	go jobs.CheckServersStart()
	go jobs.CheckPingAndTelnet()
	routes.Init()
	time.Sleep(time.Minute)
}
