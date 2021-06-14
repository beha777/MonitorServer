package settings

import "time"

type Settings struct {
	AppParams      Params           `json:"app"`
	PostgresParams PostgresSettings `json:"postgresParams"`
	BotParams      BotSettings      `json:"botParams"`
	PeriodParams   Periods          `json:"periods"`
}

type Params struct {
	ServerName    string `json:"serverName"`
	PortRun       string `json:"portRun"`
	LogFile       string `json:"logFile"`
	LogMaxSize    int    `json:"logMaxSize"`
	LogMaxBackups int    `json:"logMaxBackups"`
	LogMaxAge     int    `json:"logMaxAge"`
	LogCompress   bool   `json:"logCompress"`
}

type PostgresSettings struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	DataBase string `json:"database"`
}

type BotSettings struct {
	Url      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
	UrlID    string `json:"urlId"`
	Token    string `json:"token"`
	ChatID   string `json:"chat_id"`
}

type Periods struct {
	DefaultNotification uint      `json:"default_notification"`
	DefaultTicker       uint      `json:"default_ticker"`
	DefaultCheck        uint      `json:"default_check"`
	NilTime             time.Time `json:"nil_time"`
	UpdateServerConfig uint `json:"update_server_config"`
}
