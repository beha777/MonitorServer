package jobs

import (
	"MonitorServer/TGbot"
	"MonitorServer/client"
	"MonitorServer/db"
	"MonitorServer/models"
	"MonitorServer/server"
	"MonitorServer/settings"
	"encoding/json"
	"log"
	"net/url"
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

func CheckServices() {
	var services []models.Service
	db.GetDBConn().Find(&services, "is_active = true")
	for _, service := range services {
		if service.LastTime.Add(time.Second * time.Duration(service.CheckPeriod)).Before(time.Now()) {
			CheckServiceStatus(service)
		}
	}
}

// Check service status using 							systemctl status
func CheckServiceStatus(serviceName models.Service) {
	var curServerInfo models.ServerInfo
	db.GetDBConn().Find(&curServerInfo, "id = ?", serviceName.ServerID)
	var sendCommand client.SendCommandResponse
	url := "http://" + curServerInfo.Host + "/sendCommand?text=" + url.QueryEscape("systemctl status "+serviceName.Name+" | grep 've:'")
	err := json.Unmarshal(server.GetJson(url), &sendCommand)
	execResultString := strings.TrimSpace(sendCommand.Response)
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
				"\nStatus is " + strings.Split(execResultString, "ive: ")[1]
			serviceName.LastNotified = time.Now()
			TGbot.SendMessageToTelegramBot(message)
		}
	}
	db.GetDBConn().Save(&serviceName)
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
