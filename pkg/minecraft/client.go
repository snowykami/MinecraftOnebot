package minecraft

import (
	"MCOnebot/pkg/common"
	"errors"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
)

type ClientManager struct {
	ServerConfigs []common.ServerConfig
	AuthConfigs   map[string]common.AuthConfig
	Clients       map[string]*bot.Client // 客户端列表 name -> client
	MsgChan       chan interface{}
}

func NewClientManager(config *common.Config) *ClientManager {
	return &ClientManager{
		ServerConfigs: config.Servers,
		AuthConfigs:   config.Auth,
		Clients:       make(map[string]*bot.Client),
		MsgChan:       make(chan interface{}),
	}
}

func (c *ClientManager) JoinAllServers() {
	for _, server := range c.ServerConfigs {
		go c.JoinServer(server)
	}

}

func (c *ClientManager) JoinServer(serverConfig common.ServerConfig) {
	client := bot.NewClient()
	authConfig := c.AuthConfigs[serverConfig.Auth]
	if authConfig.Online {
		msAuth, err := GetMCcredentials(fmt.Sprintf("data/auth_%s.json", serverConfig.Auth), "88650e7e-efee-4857-b9a9-cf580a00ef43")
		if err != nil {
			return
		}
		client.Auth.UUID = msAuth.UUID
		client.Auth.AsTk = msAuth.AsTk
		client.Auth.Name = msAuth.Name
	} else {
		client.Auth.Name = authConfig.PlayerName
	}
	basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{})
	err := client.JoinServer(serverConfig.Address)
	if err != nil {
		common.Logger.Errorf("加入服务器 %s 失败: %s", serverConfig.Address, err)
	}
	common.Logger.Println("成功加入：", serverConfig.Address)
	c.Clients[serverConfig.Name] = client
	var pErr bot.PacketHandlerError
	for {
		if err = client.HandleGame(); err == nil {
			panic("处理游戏事件错误")
		}
		if errors.As(err, &pErr) {
			common.Logger.Fatalln(pErr)
		} else {
			common.Logger.Errorln(err)
		}
	}
}

func (c *ClientManager) Run() {
	c.JoinAllServers()
}
