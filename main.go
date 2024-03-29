package main

import (
	"MCOnebot/pkg/common"
	"MCOnebot/pkg/minecraft"
	"MCOnebot/pkg/onebot"
	"fmt"
	"os"
)

var (
	config = &common.Config{}
)

func main() {
	fmt.Println(
		common.Cyan(` __       __  __                                                    ______    __     
/  \     /  |/  |                                                  /      \  /  |    
$$  \   /$$ |$$/  _______    ______    _______   ______   ______  /$$$$$$  |_$$ |_   
$$$  \ /$$$ |/  |/       \  /      \  /       | /      \ /      \ $$ |_ $$// $$   |  
$$$$  /$$$$ |$$ |$$$$$$$  |/$$$$$$  |/$$$$$$$/ /$$$$$$  |$$$$$$  |$$   |   $$$$$$/   
$$ $$ $$/$$ |$$ |$$ |  $$ |$$    $$ |$$ |      $$ |  $$/ /    $$ |$$$$/      $$ | __ 
$$ |$$$/ $$ |$$ |$$ |  $$ |$$$$$$$$/ $$ \_____ $$ |     /$$$$$$$ |$$ |       $$ |/  |
$$ | $/  $$ |$$ |$$ |  $$ |$$       |$$       |$$ |     $$    $$ |$$ |       $$  $$/ 
$$/      $$/ $$/ $$/   $$/  $$$$$$$/  $$$$$$$/ $$/       $$$$$$$/ $$/         $$$$/  
  ______                       _______               __                              
 /      \                     /       \             /  |                             
/$$$$$$  | _______    ______  $$$$$$$  |  ______   _$$ |_                            
$$ |  $$ |/       \  /      \ $$ |__$$ | /      \ / $$   |                           
$$ |  $$ |$$$$$$$  |/$$$$$$  |$$    $$< /$$$$$$  |$$$$$$/                            
$$ |  $$ |$$ |  $$ |$$    $$ |$$$$$$$  |$$ |  $$ |  $$ | __                          
$$ \__$$ |$$ |  $$ |$$$$$$$$/ $$ |__$$ |$$ \__$$ |  $$ |/  |                         
$$    $$/ $$ |  $$ |$$       |$$    $$/ $$    $$/   $$  $$/                          
 $$$$$$/  $$/   $$/  $$$$$$$/ $$$$$$$/   $$$$$$/     $$$$/` + "\n\n\n"))
	err := Init()
	if err != nil {
		return
	}
	err = common.LoadConfig("config.yml", config)
	if err != nil {
		return
	}

	eventChan := make(chan common.Event, 1)
	sendChan := make(chan common.Event, 1)
	clientManager := minecraft.NewClientManager(*config)
	clientManager.EventChan = eventChan
	clientManager.SendChan = sendChan
	go clientManager.Run()

	botManager, err := onebot.NewBotManager(*config)
	botManager.EventChan = eventChan
	botManager.SendChan = sendChan
	if err != nil {
		common.Logger.Warnf("初始化 OneBot 失败: %v", err)
	}
	go botManager.Run()
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
	common.InitDatabase()
	common.Logger.Println("初始化成功")
	return nil
}
