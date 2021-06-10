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
	db.GetDBConn().Table("server").Updates(&models.Server{
		LastTime:           settings.AppSettings.PeriodParams.NilTime,
		LastNotified:       settings.AppSettings.PeriodParams.NilTime,
		NotificationPeriod: settings.AppSettings.PeriodParams.DefaultNotification,
		CheckPeriod:        settings.AppSettings.PeriodParams.DefaultCheck,
	})
	db.GetDBConn().Table("service").Updates(&models.Service{
		LastTime:           settings.AppSettings.PeriodParams.NilTime,
		LastNotified:       settings.AppSettings.PeriodParams.NilTime,
		NotificationPeriod: settings.AppSettings.PeriodParams.DefaultNotification,
		CheckPeriod:        settings.AppSettings.PeriodParams.DefaultCheck,
	})

	go jobs.CheckServicesStart()
	go jobs.CheckServersStart()
	routes.Init()
	time.Sleep(time.Minute)
}
