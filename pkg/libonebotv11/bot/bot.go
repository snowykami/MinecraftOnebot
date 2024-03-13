package bot

import (
	"github.com/sirupsen/logrus"
	"libonebotv11"
)

type Bot struct {
	// OneBot 实现名称
	Impl string
	// 机器人自身标识，为兼容 OneBot v11，此字段使用 int64 类型
	SelfID int64
	// OneBot 配置列表
	Config *libonebotv11.Config
	// 日志记录器
	Logger *logrus.Logger
}

// NewBot 创建一个新的 Bot 实例.
func NewBot(impl string, selfID int64, config *libonebotv11.Config, logger *logrus.Logger) *Bot {
	return &Bot{
		Impl:   impl,
		SelfID: selfID,
		Config: config,
		Logger: logger,
	}
}
