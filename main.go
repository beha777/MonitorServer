package main

import (
	"MonitorServer/db"
	"MonitorServer/jobs"
	"MonitorServer/models"
	"MonitorServer/routes"
	"MonitorServer/settings"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

func main() {
	settings.AppSettings = settings.ReadSettings()
	logSettings()

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

func logSettings() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   settings.AppSettings.AppParams.LogFile,
		MaxSize:    settings.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: settings.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     settings.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   settings.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	})
}
