package server

import (
	"MonitorServer/TGbot"
	"MonitorServer/client"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func GetJson(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	log.Println(url)
	if err != nil {
		message := "❌ Server doens't respond: " +
			"\n" + err.Error()
		TGbot.SendMessageToTelegramBot(message)
		log.Println("resp_DO_err", err)
		return []byte{}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		message := "❌ Server doens't respond: " +
			"\n" + err.Error()
		TGbot.SendMessageToTelegramBot(message)
		log.Fatalln("resp_ReadAll_err", err)
		return []byte{}
	}
	defer resp.Body.Close()
	return b
}

func GetCPUload(host string) float64 {
	var sendCommand client.SendCommandResponse
	url := "http://" + host + "/sendCommand?text=" + url.QueryEscape("top -bn 1 | fgrep 'load'")
	err := json.Unmarshal(GetJson(url), &sendCommand)
	if err != nil {
		log.Println("EXEC_GetCPUload error", err)
	} else {
		execResultString := strings.Split(sendCommand.Response, "average: ")[1]
		log.Println("CPU_execResultString =", execResultString)
		log.Println(execResultString)
		re := regexp.MustCompile(`[-]?\d+[.,]?\d*`)
		parsedValues := re.FindAllString(execResultString, -1)
		curCPUload, err := strconv.ParseFloat(strings.Replace(parsedValues[1], ",", ".", -1), 64)
		log.Println("curCPUload =", curCPUload)
		if err != nil {
			log.Println("CPU_CONV error", err)
		} else {
			log.Printf("%.2f%%\n", curCPUload*100)
			return curCPUload * 100
		}
	}
	return -1
}

func GetMemLoad(host string) float64 {
	var sendCommand client.SendCommandResponse
	url := "http://" + host + "/sendCommand?text=" + url.QueryEscape("top -bn 1 | fgrep 'Mem :'")
	err := json.Unmarshal(GetJson(url), &sendCommand)
	if err != nil {
		log.Println("EXEC_GetMemLoad error", err)
	} else {
		execResultString := sendCommand.Response
		log.Println(execResultString)
		re := regexp.MustCompile(`[-]?\d+[.,]?\d*`)
		parsedValues := re.FindAllString(execResultString, -1)
		curMemUsed, err := strconv.ParseFloat(parsedValues[2], 64)
		if err != nil {
			log.Println("Mem_CONV_1 error", err)
		}
		MemTotal, err := strconv.ParseFloat(parsedValues[0], 64)
		if err != nil {
			log.Println("Mem_CONV_2 error", err)
		} else {
			curMemLoad := curMemUsed / MemTotal * 100
			log.Printf("%.2f%%\n", curMemLoad)
			return curMemLoad
		}
	}
	return -1
}
func GetDiscUsage(host string) float64 {
	var sendCommand client.SendCommandResponse
	url := "http://" + host + "/sendCommand?text=" + url.QueryEscape("df -h /")
	err := json.Unmarshal(GetJson(url), &sendCommand)
	if err != nil {
		log.Println("EXEC_GetDiscUsage error", err)
	} else {
		execResultString := sendCommand.Response
		log.Println(execResultString)
		re := regexp.MustCompile(`[-]?\d+[.,]?\d*`)
		parsedValues := re.FindAllString(execResultString, -1)
		curDiscUsed, err := strconv.ParseFloat(parsedValues[3], 64)
		if err != nil {
			log.Println("Disc_CONV error", err)
		} else {
			log.Printf("%.2f%%\n", curDiscUsed)
			return curDiscUsed
		}
	}
	return -1
}
