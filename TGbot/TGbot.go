package TGbot

import (
	"MonitorServer/settings"
	"log"
	"net/http"
	"net/url"
)

func SendMessageToTelegramBot(message string) {
	message = url.QueryEscape(message)
	_, err := http.Get("https://api.telegram.org/bot" + settings.AppSettings.BotParams.Token + "/sendMessage?chat_id=" + settings.AppSettings.BotParams.ChatID + "&text=" + message)
	if err != nil {
		log.Println("SendMessageToTelegramBot(", message, ") error", err)
	}

}
