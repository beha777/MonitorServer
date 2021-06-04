package jobs

import (
	"MonitorServer/db"
	"MonitorServer/models"
	"MonitorServer/server"
	"github.com/reiver/go-telnet"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CheckServicesStart() {
	timer := time.NewTicker(time.Second * 10)
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
	timer := time.NewTicker(time.Second * 10)
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

func CheckPingAndTelnet () {
	timer := time.NewTicker(time.Second * 10)
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
	db.GetDBConn().Find(&servers)
	for _, curServer := range servers {
		out, err := exec.Command("ping", strings.Split(curServer.Host, ":")[0]).Output()
		if err != nil && !strings.Contains(string(out), "Lost = 0") {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20Cant%20ping%3A%20" + strings.Split(curServer.Host, ":")[0])
			if err != nil {
				log.Println("GET_ping error", err)
			}

			log.Println("PING error", err)
		}
	}
}

func CheckTelnet() {
	var servers []models.ServerInfo
	db.GetDBConn().Find(&servers)
	for _, curServer := range servers {
		_, err := telnet.DialTo(curServer.Host)
		if err != nil {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20Cant%20telnet%20port%3A%20" + strings.Split(curServer.Host, ":")[1])
			if err != nil {
				log.Println("GET_telnet error", err)
			}
			log.Println("TELNET error", err)
		}
	}
}



func CheckServers() {
	var servers []models.Server
	db.GetDBConn().Find(&servers)
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
		_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20Cant%20parse%20" + curServer.Param)
		if err != nil {
			log.Println("GET_", curServer.Param, " parse error", err)
		}
	}
	curServer.LastTime = time.Now()
	if curServer.LastNotified.Add(time.Second * time.Duration(curServer.NotificationPeriod)).Before(time.Now()) {
		curServer.LastNotified = time.Now()
		if (curServer.Condition == ">" && curParam > curServer.Limit) || (curServer.Condition == "<" && curParam < curServer.Limit) {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20ServerID%3A" + strconv.Itoa(int(curServer.ServerID)) + "%3A" + curServer.Param + curServer.Condition + strconv.Itoa(int(curServer.Limit)))
			if err != nil {
				log.Println("GET_1 error", err)
			}
		} else {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9C%94%EF%B8%8F" + "%20ServerID%3A" + strconv.Itoa(int(curServer.ServerID)) + "%3A" + curServer.Param + curServer.Condition + strconv.Itoa(int(curServer.Limit)))
			if err != nil {
				log.Println("GET_2 error", err)
			}
		}
	}
	db.GetDBConn().Save(&curServer)
}

func CheckServices() {
	var services []models.Service
	db.GetDBConn().Find(&services)
	for _, service := range services {
		if service.LastTime.Add(time.Second * time.Duration(service.CheckPeriod)).Before(time.Now()) {
			CheckServiceStatus(service)
		}
	}
}

func CheckServiceStatus(serviceName models.Service) {
	curServerCon := server.ConnectToServer(serviceName.ServerID)
	execResult, err := curServerCon.Exec("systemctl is-active " + serviceName.Name)
	execResultString := string(execResult)
	serviceName.LastTime = time.Now()

	if err != nil && !strings.Contains(execResultString, "active") {
		log.Println("EXEC_CheckServiceStatus error-", err, execResultString)
	}
	if serviceName.LastNotified.Add(time.Second * time.Duration(serviceName.NotificationPeriod)).Before(time.Now()) {
		serviceName.LastNotified = time.Now()
		if execResultString != serviceName.State {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20ServerID%3A" + strconv.Itoa(int(serviceName.ServerID)) + "%3A" + serviceName.Name + "%20Status%20is%20" + execResultString)
			if err != nil {
				log.Println("GET_3 error", err)
			}
		} else {
			_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9C%94%EF%B8%8F" + "%20ServerID%3A" + strconv.Itoa(int(serviceName.ServerID)) + "%3A" + serviceName.Name + "%20Status%20is%20" + execResultString)
			if err != nil {
				log.Println("GET_4 error", err)
			}
		}
	}
	db.GetDBConn().Save(&serviceName)
}
