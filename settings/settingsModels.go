package settings

type Settings struct {
	AppParams      Params           `json:"app"`
	PostgresParams PostgresSettings `json:"postgresParams"`
	BotParams      BotSettings      `json:"botParams"`
	PeriodParams   Periods          `json:"periods"`
}

type Params struct {
	ServerName string `json:"serverName"`
	PortRun    string `json:"portRun"`
}

type PostgresSettings struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	DataBase string `json:"database"`
}

type BotSettings struct {
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
}

type Periods struct {
	DefaultNotification int `json:"default_notification"`
	DefaultTicker       int `json:"default_ticker"`
	DefaultCheck        int `json:"default_check"`
}
