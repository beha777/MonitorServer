package models

import "time"

type Server struct {
	ID                 uint `gorm:"primaryKey"`
	ServerID           uint
	Param              string
	Condition          string
	Limit              float64
	LastTime           time.Time
	CheckPeriod        uint
	LastNotified       time.Time
	NotificationPeriod uint
	IsActive           bool
}
type Service struct {
	ID                 uint `gorm:"primaryKey"`
	ServerID           uint
	Name               string
	State              string
	LastTime           time.Time
	CheckPeriod        uint
	LastNotified       time.Time
	NotificationPeriod uint
	IsActive           bool
}

type ServerInfo struct {
	ID         uint    `json:"id"`
	Host       string  `json:"host"`
	Version    float64 `json:"version"`
	Owners     string  `json:"owners"`
	IsActive   bool    `json:"is_active"`
	TgUrl      string  `json:"tg_url"`
	TgLogin    string  `json:"tg_login"`
	TgPassword string  `json:"tg_password"`
	TgUrlId    string  `json:"tg_url_id"`
	TgBotToken string  `json:"tg_bot_token"`
	TgChatId   string  `json:"tg_chat_id"`
	Log        bool    `json:"log"`
}
