package server

import (
	"MonitorServer/db"
	"MonitorServer/models"
	"github.com/sfreiberg/simplessh"
	"log"
	"net/http"
)

var centosServer *simplessh.Client

func InitServer() *simplessh.Client {
	server, err := simplessh.ConnectWithPassword("127.0.0.1:2281", "root", "q1w2r3t4")
	if err != nil {
		log.Println("CONNECT error", err)
	}
	//defer centosServer.Close()
	return server
}

func ConnectToServer(serverID uint) *simplessh.Client {
	var curServer models.ServerInfo
	db.GetDBConn().Find(&curServer, serverID)

	server, err := simplessh.ConnectWithPassword(curServer.Host, curServer.Login, "q1w2r3t4")
	if err != nil {
		_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20Cant%20connect%20server%3A%20" + curServer.Host + "%3A" + curServer.Login)
		if err != nil {
			log.Println("GET_connect error", err)
		}
		log.Println("CONNECT_error", err)
	}
	return server
}

func ConnectServer() {
	centosServer = InitServer()
}

func GetServerConn() *simplessh.Client {
	return centosServer
}
