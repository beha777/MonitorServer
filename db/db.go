package db

import (
	"MonitorServer/Cipher"
	"MonitorServer/models"
	"MonitorServer/settings"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
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
	db, err := gorm.Open("postgres", connString)

	if err != nil {
		_, err := http.Get("https://api.telegram.org/bot1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8/sendMessage?chat_id=@ServerParamStatus&text=%E2%9D%8C" + "%20Cant%20connect%20server%3A%20" + settingParams.Server + "%3A" + settingParams.DataBase)
		if err != nil {
			log.Println("GET_DBconnect error", err)
		}
		log.Fatal("Couldn't connect to postgresql database", err.Error(), settingParams.Server)
		//logger.Error.Println()
	}
	db.LogMode(true)
	db.SingularTable(true)
	return db
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

func AddServices(serviceList []string, serverId uint) {
	for _, serviceName := range serviceList {
		newService := models.Service{
			ServerID:           serverId,
			Name:               strings.Split(strings.TrimSpace(serviceName), " ")[0],
			State:              "active",
			LastTime:           time.Time{},
			CheckPeriod:        3,
			LastNotified:       time.Time{},
			NotificationPeriod: 3600,
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
