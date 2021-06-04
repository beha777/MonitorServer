package server

import (
	"github.com/sfreiberg/simplessh"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func GetCPUload(centosServer *simplessh.Client) float64 {
	execResult, err := centosServer.Exec("top -bn 1 | fgrep 'load'")
	execResultString := string(execResult)
	if err != nil {
		log.Println("EXEC_GetCPUload error", err)
	} else {
		log.Println(execResultString)
		re := regexp.MustCompile(`[-]?\d+[.,]?\d*`)
		parsedValues := re.FindAllString(execResultString, -1)
		curCPUload, err := strconv.ParseFloat(strings.Replace(parsedValues[6], ",", ".", -1), 64)
		if err != nil {
			log.Println("CPU_CONV error", err)
		} else {
			log.Printf("%.2f%%\n", curCPUload*100)
			return curCPUload * 100
		}
	}
	return -1
}

func GetMemLoad(centosServer *simplessh.Client) float64 {
	execResult, err := centosServer.Exec("top -bn 1 | fgrep 'Mem :'")
	execResultString := string(execResult)
	if err != nil {
		log.Println("EXEC_GetMemLoad error", err)
	} else {
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
func GetDiscUsage(centosServer *simplessh.Client) float64 {
	execResult, err := centosServer.Exec("df -h /")
	execResultString := string(execResult)
	if err != nil {
		log.Println("EXEC_GetDiscUsage error", err)
	} else {
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

func GetServicesList(centosServer *simplessh.Client) []string {
	execResult, err := centosServer.Exec("systemctl list-units --type=service --all")
	execResultString := string(execResult)
	var execResultStringTrimmed []string
	if err != nil {
		log.Println("EXEC_GetServicesList error", err)
	} else {
		execResultString = regexp.MustCompile(`[ ●]+`).ReplaceAllString(execResultString, " ")
		execResultStringTrimmed = strings.Split(execResultString, "\n")
		log.Println(execResultStringTrimmed[1])
		log.Println(strings.Split(strings.TrimSpace(execResultStringTrimmed[2]), " ")[1])

	}
	return execResultStringTrimmed
}
