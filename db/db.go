package db

import (
	"MonitorServer/Cipher"
	"MonitorServer/TGbot"
	"MonitorServer/models"
	"MonitorServer/settings"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"strings"
	"time"
)

var database *gorm.DB

func initDB() *gorm.DB {
	settingParams := settings.AppSettings.PostgresParams
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		settingParams.Server, settingParams.Port,
		settingParams.User, settingParams.DataBase,
		settingParams.Password)
	dbase, err := gorm.Open("postgres", connString)

	if err != nil {
		message := "❌ Can't connect database server: " + settingParams.Server + "\nDatabase: " + settingParams.DataBase
		TGbot.SendMessageToTelegramBot(message)
		log.Fatal("Couldn't connect to postgresql database", err.Error(), settingParams.Server)
	}
	dbase.LogMode(true)
	dbase.SingularTable(true)
	return dbase
}

func ConnectDatabase() {
	database = initDB()
}

func GetDBConn() *gorm.DB {
	return database
}

func AddServer(NewServer models.ServerInfo) bool {
	if database.Find(&models.ServerInfo{}, "host = ? and login = ?", NewServer.Host, NewServer.Login).Error == nil {
		log.Println("SERVER exists")
		return false
	}
	NewServer.Password = Cipher.Encode(NewServer.Password)
	database.Create(&NewServer)
	return true
}

var nilTime = time.Unix(1, 0)

func AddServices(serviceList []string, serverId uint) {
	for _, serviceName := range serviceList {
		newService := models.Service{
			ServerID:           serverId,
			Name:               strings.Split(strings.TrimSpace(serviceName), " ")[0],
			State:              "active",
			LastTime:           nilTime,
			CheckPeriod:        uint(settings.AppSettings.PeriodParams.DefaultCheck),
			LastNotified:       nilTime,
			NotificationPeriod: uint(settings.AppSettings.PeriodParams.DefaultNotification),
		}
		if newService.Name == "UNIT" || database.Find(&models.Service{}, "Name = ?", newService.Name).Error == nil {
			continue
		}
		if len(newService.Name) < 2 {
			break
		}
		database.Create(&newService)
	}
}
