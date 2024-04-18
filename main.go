package main

import (
	"MCOnebot/pkg/common"
	"MCOnebot/pkg/minecraft"
)

func main() {
	common.Logger.Info("Starting MCOnebot...")
	common.Logger.Info("Loading config...")
	config, err := common.LoadConfig()
	if err != nil {
		common.Logger.Error("Failed to load config: ", err)
		return
	}
	common.Logger.Infof("Config loaded: %+v", config)

	cm := minecraft.NewClientManager(config.Minecraft)
	go cm.Start()
	select {}
}
