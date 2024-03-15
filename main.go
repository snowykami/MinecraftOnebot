package main

import (
	"MCOnebot/pkg/common"
	mc "MCOnebot/pkg/minecraft"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	config = &common.Config{}
)

func main() {
	err := Init()
	if err != nil {
		return
	}
	err = common.LoadConfig("config.yml", config)
	if err != nil {
		return
	}

	clientManager := mc.NewClientManager(*config)
	go clientManager.Run()

	select {}
}

func Init() error {
	// 初始化检测并创建必要的文件夹，有则跳过
	folders := []string{"data", "logs"}
	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			err := os.Mkdir(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	common.Logger.Println("初始化成功")
	return nil
}

// ConnectOnebot 连接 Onebot
func ConnectOnebot() {
	common.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	common.Logger.SetOutput(os.Stdout)
	common.Logger.SetLevel(logrus.DebugLevel)
}
