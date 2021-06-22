package jobs

import (
	"MonitorServer/client"
	"MonitorServer/db"
	"MonitorServer/models"
	"MonitorServer/server"
	"MonitorServer/settings"
	"encoding/json"
	"log"
	"time"
)

func UpdateClientConfigStart() {
	timer := time.NewTicker(time.Second * time.Duration(settings.AppSettings.PeriodParams.UpdateClientConfig))
	defer timer.Stop()

	log.Println("UpdateClientConfigStartED")
	for {
		select {
		case <-timer.C:
			log.Println("UpdateClientConfig")
			var servers []models.ServerInfo
			db.GetDBConn().Find(&servers)
			for _, curServer := range servers {
				var curClient client.Config
				url := "http://" + curServer.Host + "/getParams"
				err := json.Unmarshal(server.GetJson(url), &curClient)
				if err != nil {
					log.Println("/getParams error", err)
				} else if curClient.AppParams.Version > curServer.Version {
					log.Printf("--------%f\n%f", curClient.AppParams.Version, curServer.Version)
					UpdateClientConfig(curServer, curClient)
				}
			}
		}
	}
}

func UpdateClientConfig(curServer models.ServerInfo, curClient client.Config) {
	curServer = models.ServerInfo{
		ID:         curServer.ID,
		Host:       curServer.Host,
		Version:    curClient.AppParams.Version,
		Owners:     curClient.BotParams.Owners,
		IsActive:   curClient.AppParams.Active,
		TgUrl:      curClient.BotParams.Url,
		TgLogin:    curClient.BotParams.Login,
		TgPassword: curClient.BotParams.Password,
		TgUrlId:    curClient.BotParams.UrlId,
		TgBotToken: curClient.BotParams.Token,
		TgChatId:   curClient.BotParams.ChatID,
		Log:        curClient.AppParams.Log,
	}
	db.GetDBConn().Save(&curServer)

	//reCreate all services
	db.GetDBConn().Delete(&models.Service{}, "server_id = ?", curServer.ID)
	for _, service := range curClient.Services {
		log.Printf("service: %+v", service)
		db.GetDBConn().Create(&models.Service{
			ServerID:           curServer.ID,
			Name:               service.Name,
			State:              service.State,
			LastTime:           settings.AppSettings.PeriodParams.NilTime,
			CheckPeriod:        service.CheckPeriod,
			LastNotified:       settings.AppSettings.PeriodParams.NilTime,
			NotificationPeriod: service.NotificationPeriod,
			IsActive:           service.Active,
		})
	}

	//reCreate all serverParams
	db.GetDBConn().Delete(&models.Server{}, "server_id = ?", curServer.ID)
	for _, serverParam := range curClient.ServerParams {

		log.Printf("param: %+v", serverParam)
		db.GetDBConn().Create(&models.Server{
			ServerID:           curServer.ID,
			Param:              serverParam.Name,
			Condition:          serverParam.Condition,
			Limit:              serverParam.Limit,
			LastTime:           settings.AppSettings.PeriodParams.NilTime,
			CheckPeriod:        serverParam.CheckPeriod,
			LastNotified:       settings.AppSettings.PeriodParams.NilTime,
			NotificationPeriod: serverParam.NotificationPeriod,
			IsActive:           serverParam.Active,
		})
	}
	for _, host := range curClient.Hosts {
		db.GetDBConn().Create(&models.Server{
			ServerID:           curServer.ID,
			Param:              host.Name,
			Condition:          "OK",
			LastTime:           settings.AppSettings.PeriodParams.NilTime,
			CheckPeriod:        host.CheckPeriod,
			LastNotified:       settings.AppSettings.PeriodParams.NilTime,
			NotificationPeriod: host.NotificationPeriod,
			IsActive:           host.Active,
		})
	}
}
