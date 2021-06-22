package server

import (
	"MonitorServer/TGbot"
	"github.com/sfreiberg/simplessh"
	"log"
)

var centosServer *simplessh.Client

func InitServer() *simplessh.Client {
	server, err := simplessh.ConnectWithPassword("127.0.0.1:2281", "root", "q1")
	if err != nil {
		message := "❌ Can't connect server: 127.0.0.1:2281\nUsername: root"
		TGbot.SendMessageToTelegramBot(message)
		log.Println("CONNECT error", err)
	}
	//defer centosServer.Close()
	return server
}

/*
func ConnectToServer(serverID uint) *simplessh.Client {
	var curServer models.ServerInfo
	db.GetDBConn().Find(&curServer, serverID)
	server, err := simplessh.ConnectWithPassword(curServer.Host, curServer.Login, Cipher.Decode(curServer.Password))
	if err != nil {
		message := "❌ Can't connect server: " + curServer.Host + "\nUsername: " + curServer.Login +
			"\n" + curServer.Owners
		TGbot.SendMessageToTelegramBot(message)
		log.Println("CONNECT_error", err)
	}
	return server
}*/

func ConnectServer() {
	centosServer = InitServer()
}

func GetServerConn() *simplessh.Client {
	return centosServer
}
