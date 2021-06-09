package TGbot

import (
	"MonitorServer/settings"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func SendMessageToTelegramBot(message string) {
	//newline := "%0A"
	url4tb := settings.AppSettings.BotParams.Url

	req, err := http.NewRequest("GET", url4tb, nil)
	if err != nil {
		fmt.Println("1SendMessageToTelegramBot(", message, ") error", err)
	}

	req.Header.Add("Login", settings.AppSettings.BotParams.Login)
	req.Header.Add("Password", settings.AppSettings.BotParams.Password)
	req.Header.Add("UrlID", settings.AppSettings.BotParams.UrlID)
	req.Header.Add("Msg", url.QueryEscape(message))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("2SendMessageToTelegramBot(", message, ") error", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("3SendMessageToTelegramBot(", message, ") error", err)
	}

	log.Println("SendMessageToTelegramBot(", message, "):", string(body))

	/*

		message = url.QueryEscape(message)
		_, err := http.Get("https://api.telegram.org/bot" + settings.AppSettings.BotParams.Token + "/sendMessage?chat_id=" + settings.AppSettings.BotParams.ChatID + "&text=" + message)
		if err != nil {
			log.Println("SendMessageToTelegramBot(", message, ") error", err)
		}
	*/
}
