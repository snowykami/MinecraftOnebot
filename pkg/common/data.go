package common

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var PlayerDB *gorm.DB
var MessageDB *gorm.DB

func InitDatabase() {
	PlayerDB, err := gorm.Open(sqlite.Open("data/players.db"), &gorm.Config{})
	err = PlayerDB.AutoMigrate(Player{})
	if err != nil {
		Logger.Errorf("数据库迁移失败: %s", err)
	}
	MessageDB, err := gorm.Open(sqlite.Open("data/messages.db"), &gorm.Config{})
	err = MessageDB.AutoMigrate(Message{})
}
