package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

var AppSettings Settings

// ReadSettings to init app settings
func ReadSettings() Settings {
	var appParams Settings
	doc, err := ioutil.ReadFile("./settings-dev.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = json.Unmarshal(doc, &appParams)
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
	appParams.PeriodParams.NilTime = time.Unix(1, 0)
	return appParams
}
