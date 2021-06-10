package jobs

import (
	"MonitorServer/TGbot"
	"MonitorServer/db"
	"MonitorServer/models"
	"MonitorServer/server"
	"MonitorServer/settings"
	"github.com/reiver/go-telnet"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CheckServicesStart() {
	timer := time.NewTicker(time.Second * time.Duration(settings.AppSettings.PeriodParams.DefaultTicker))
	defer timer.Stop()

	log.Println("CheckServicesStartED")
	for {
		select {
		case <-timer.C:
			CheckServices()
		}
	}
}

func CheckServersStart() {
	timer := time.NewTicker(time.Second * time.Duration(settings.AppSettings.PeriodParams.DefaultTicker))
	defer timer.Stop()

	log.Println("CheckServersStartED")
	for {
		select {
		case <-timer.C:
			log.Println("CheckServers")
			CheckServers()
		}
	}
}

func CheckPingAndTelnet() {
	timer := time.NewTicker(time.Second * time.Duration(settings.AppSettings.PeriodParams.DefaultTicker))
	defer timer.Stop()

	log.Println("CheckPingAndTelnetED")
	for {
		select {
		case <-timer.C:
			CheckPing()
			CheckTelnet()
		}
	}
}

func CheckPing() {
	var servers []models.ServerInfo
	db.GetDBConn().Find(&servers, "is_active = true")
	for _, curServer := range servers {
		out, err := exec.Command("ping", strings.Split(curServer.Host, ":")[0]).Output()
		if err != nil && !strings.Contains(string(out), "Lost = 0") {
			message := "❌ Can't ping: " + strings.Split(curServer.Host, ":")[0] +
				"\n" + curServer.Owners
			TGbot.SendMessageToTelegramBot(message)
			log.Println("PING error", err)
		}
	}
}

func CheckTelnet() {
	var servers []models.ServerInfo
	db.GetDBConn().Find(&servers, "is_active = true")
	for _, curServer := range servers {
		_, err := telnet.DialTo(curServer.Host)
		if err != nil {
			message := "❌ Can't telnet port: " + strings.Split(curServer.Host, ":")[1] +
				"\n" + curServer.Owners
			TGbot.SendMessageToTelegramBot(message)
			log.Println("TELNET error", err)
		}
	}
}

func CheckServers() {
	var servers []models.Server
	var activeServerIds []uint
	db.GetDBConn().Table("server_info").Where("is_active = false").Pluck("id", &activeServerIds)
	if len(activeServerIds) == 0 {
		activeServerIds = []uint{0}
	}
	db.GetDBConn().Not("server_id", activeServerIds).Find(&servers)
	for _, curServer := range servers {
		if curServer.LastTime.Add(time.Second * time.Duration(curServer.CheckPeriod)).Before(time.Now()) {
			CheckServerStatus(curServer)
		}
	}
}

func CheckServerStatus(curServer models.Server) {
	var curParam float64
	curServerCon := server.ConnectToServer(curServer.ServerID)
	switch {
	case curServer.Param == "CpuLoad":
		curParam = server.GetCPUload(curServerCon)
	case curServer.Param == "MemLoad":
		curParam = server.GetMemLoad(curServerCon)
	case curServer.Param == "DiscUsed":
		curParam = server.GetDiscUsage(curServerCon)
	}
	log.Println("curParam ---", curParam)
	if curParam == -1 {
		message := "❌ Can't parse " + curServer.Param
		TGbot.SendMessageToTelegramBot(message)
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

func CheckServices() {
	var services []models.Service
	db.GetDBConn().Find(&services, "is_active = true")
	for _, service := range services {
		if service.LastTime.Add(time.Second*time.Duration(service.CheckPeriod)).Before(time.Now()) {
			CheckServiceStatus(service)
		}
	}
}
// Check service status using 							systemctl is-active
/*func CheckServiceStatus(serviceName models.Service) {
	curServerCon := server.ConnectToServer(serviceName.ServerID)
	execResult, err := curServerCon.Exec("systemctl is-active " + serviceName.Name)
	execResultString := strings.TrimSpace(string(execResult))
	serviceName.LastTime = time.Now()
	if err != nil && !strings.Contains(execResultString, "active") {
		log.Println("EXEC_CheckServiceStatus error-", err, execResultString)
	}
	if serviceName.LastNotified.Add(time.Second * time.Duration(serviceName.NotificationPeriod)).Before(time.Now()) {
		if execResultString != serviceName.State {
			var curServerInfo models.ServerInfo
			db.GetDBConn().Find(&curServerInfo, "id = ?", serviceName.ServerID)
			message := "❌ ServerID:" + strconv.Itoa(int(serviceName.ServerID)) +
				"\nHost: " + curServerInfo.Host +
				"\nService: " + serviceName.Name +
				"\nStatus is " + execResultString
			serviceName.LastNotified = time.Now()
			TGbot.SendMessageToTelegramBot(message)
		}
	}
	db.GetDBConn().Save(&serviceName)
}*/
// Check service status using 							systemctl status
func CheckServiceStatus(serviceName models.Service) {
	curServerCon := server.ConnectToServer(serviceName.ServerID)
	execResult, err := curServerCon.Exec("systemctl status " + serviceName.Name + " | grep 've:'")
	execResultString := strings.TrimSpace(string(execResult))
	var serviceStatus string
	if err != nil && !strings.Contains(execResultString, "could not be found.") {
		log.Println("EXEC_CheckServiceStatus error-", err, execResultString)
	} else {
		serviceStatus = strings.Split(execResultString, "ive: ")[1]
		serviceStatus = strings.Split(serviceStatus, " ")[0]
	}
	serviceName.LastTime = time.Now()
	if serviceName.LastNotified.Add(time.Second * time.Duration(serviceName.NotificationPeriod)).Before(time.Now()) {
		if serviceStatus != serviceName.State {
			var curServerInfo models.ServerInfo
			db.GetDBConn().Find(&curServerInfo, "id = ?", serviceName.ServerID)
			message := "❌ ServerID:" + strconv.Itoa(int(serviceName.ServerID)) +
				"\nHost: " + curServerInfo.Host +
				"\nService: " + serviceName.Name +
				"\nStatus is " + strings.Split(execResultString, "ive: ")[1] +
				"\n" + serviceName.Owners
			serviceName.LastNotified = time.Now()
			TGbot.SendMessageToTelegramBot(message)
		}
	}
	db.GetDBConn().Save(&serviceName)
}