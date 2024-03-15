package main

import (
	"MCOnebot/pkg/common"
	mc "MCOnebot/pkg/minecraft"
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
