package common

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name  string `gorm:"primaryKey"`
	UUID  string `gorm:"uniqueIndex"`
	IntId int64  `gorm:"primaryKey;autoIncrement"`
}

type Message struct {
	gorm.Model
	MessageID  int64  `gorm:"primaryKey"`
	PlayerName string `json:"player_name"`
	ServerName string `json:"server_name"`
	Context    string `json:"context"`
	Type       string `json:"type"`
}
