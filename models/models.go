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
}

type ServerInfo struct {
	ID       uint   `json:"id":"primaryKey"`
	Host     string `json:"host"`
	Login    string `json:"login"`
	Password string `json:"password"`
	OS       string `json:"os"`
	Version  string `json:"version"`
}
