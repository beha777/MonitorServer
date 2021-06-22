package jobs

import (
	"MonitorServer/TGbot"
	"MonitorServer/db"
	"MonitorServer/models"
	"MonitorServer/server"
	"MonitorServer/settings"
	"log"
	"strconv"
	"strings"
	"time"
)

func CheckServerParamsStart() {
	timer := time.NewTicker(time.Second * time.Duration(settings.AppSettings.PeriodParams.DefaultTicker))
	defer timer.Stop()

	log.Println("CheckServersStartED")
	for {
		select {
		case <-timer.C:
			log.Println("CheckServers")
			CheckServerParams()
		}
	}
}

func CheckServerParams() {
	var servers []models.Server
	var activeServerIds []uint
	db.GetDBConn().Table("server_info").Where("is_active = false").Pluck("id", &activeServerIds)
	if len(activeServerIds) == 0 {
		activeServerIds = []uint{0}
	}
	db.GetDBConn().Not("server_id", activeServerIds).Find(&servers)
	for _, curServer := range servers {
		if curServer.LastTime.Add(time.Second * time.Duration(curServer.CheckPeriod)).Before(time.Now()) {
			CheckServerParamStatus(curServer)
		}
	}
}

func CheckServerParamStatus(curServer models.Server) {
	var curParam float64
	var curServerInfo models.ServerInfo
	db.GetDBConn().Find(&curServerInfo, "id = ?", curServer.ServerID)
	switch {
	case curServer.Param == "CpuLoad":
		curParam = server.GetCPUload(curServerInfo.Host)
	case curServer.Param == "MemLoad":
		curParam = server.GetMemLoad(curServerInfo.Host)
	case curServer.Param == "DiscUsed":
		curParam = server.GetDiscUsage(curServerInfo.Host)
	}
	log.Println("curParam ---", curParam)
	if curParam == -1 {
		hostPort := strings.Split(curServer.Param, ":")[1]
		if err := CheckTelnet(hostPort); err != nil {
			message := "❌ Can't telnet port: " + hostPort +
				"\n" + curServerInfo.Owners
			TGbot.SendMessageToTelegramBot(message)
			log.Println("TELNET error", err)
		}
		hostIp := strings.Split(curServer.Param, ":")[0]
		if err := CheckPing(hostIp); err != nil {
			message := "❌ Can't ping: " + hostIp +
				"\n" + curServerInfo.Owners
			TGbot.SendMessageToTelegramBot(message)
			log.Println("PING error", err)
		}

	}
	curServer.LastTime = time.Now()
	if curServer.LastNotified.Add(time.Second * time.Duration(curServer.NotificationPeriod)).Before(time.Now()) {
		curServer.LastNotified = time.Now()
		if (curServer.Condition == ">" && curParam > curServer.Limit) || (curServer.Condition == "<" && curParam < curServer.Limit) {
			curServer.LastNotified = time.Now()
			var curServerInfo models.ServerInfo
			db.GetDBConn().Find(&curServerInfo, "id = ?", curServer.ServerID)
			message := "❌ ServerID:" + strconv.Itoa(int(curServer.ServerID)) +
				"\nHost: " + curServerInfo.Host +
				"\nParam:" + curServer.Param + " = " + strconv.FormatFloat(curParam, 'f', 1, 64) + " ( " + curServer.Condition + strconv.Itoa(int(curServer.Limit)) +
				")\n" + curServerInfo.Owners
			TGbot.SendMessageToTelegramBot(message)
		}
	}
	db.GetDBConn().Save(&curServer)
}
