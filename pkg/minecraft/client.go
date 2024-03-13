package minecraft

import (
	"MCOnebot/pkg/common"
	"MCOnebot/pkg/config"
	"errors"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"log"
)

type ClientManager struct {
	ServerConfigs []config.ServerConfig
	Clients       map[string]*bot.Client // 客户端列表 name -> client
	MsgChan       chan interface{}
}

func (c *ClientManager) JoinServer(serverConfig config.ServerConfig) {
	client := bot.NewClient()
	basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{})
	err := client.JoinServer(serverConfig.Address)
	if err != nil {
		common.Logger.Errorf("加入服务器 %s 失败: %s", serverConfig.Address, err)
	}
	log.Println("成功加入：", serverConfig.Address)
	c.Clients[serverConfig.Name] = client
	var pErr bot.PacketHandlerError
	for {
		if err = client.HandleGame(); err == nil {
			panic("处理游戏事件错误")
		}
		if errors.As(err, &pErr) {
			log.Print(pErr)
		} else {
			log.Fatal(err)
		}
	}
}

func (c *ClientManager) Run() {
	for _, server := range c.ServerConfigs {
		go c.JoinServer(server)
	}
}
