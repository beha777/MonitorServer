package main

import (
	"MonitorServer/db"
	"MonitorServer/jobs"
	"MonitorServer/models"
	"MonitorServer/routes"
	"MonitorServer/server"
	"MonitorServer/settings"
	"log"
	"time"
)

func main() {

	settings.AppSettings = settings.ReadSettings()
	server.ConnectServer()
	db.ConnectDatabase()
	//db.GetDBConn().DropTable(&models.Server{}, &models.ServerInfo{}, &models.Service{})
	db.GetDBConn().AutoMigrate(&models.Server{}, &models.ServerInfo{}, &models.Service{})
	routes.Init()
	var centOsServer = models.ServerInfo{
		Host:     "127.0.0.1:2281",
		Login:    "root",
		Password: "q1w2r3t4",
		OS:       "CentOS",
		Version:  "7.0",
	}

	db.AddServer(centOsServer)
	curCPUload := server.GetCPUload(server.GetServerConn())
	curMemLoad := server.GetMemLoad(server.GetServerConn())
	curDiscUsed := server.GetDiscUsage(server.GetServerConn())

	var serversInfo []models.ServerInfo
	db.GetDBConn().Find(&serversInfo)
	for _, serverInfo := range serversInfo {
		db.AddServices(server.GetServicesList(server.ConnectToServer(serverInfo.ID)), serverInfo.ID)
	}

	go jobs.CheckServicesStart()
	go jobs.CheckServersStart()
	go jobs.CheckPingAndTelnet()

	log.Printf("CPU_load = %.0f%%\n"+
		"Mem_load = %.0f%%\n"+
		"Disc_load = %.0f%%\n",
		curCPUload, curMemLoad, curDiscUsed)

	time.Sleep(time.Minute)
}
