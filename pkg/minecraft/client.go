package minecraft

import (
	"MCOnebot/pkg/common"
	"errors"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"time"
)

var defaultName = "MCOnebot"

type ClientManager struct {
	ServerConfigs []common.ServerConfig
	AuthConfigs   map[string]common.AuthConfig
	Clients       map[string]*bot.Client // 客户端列表 name -> client
	MsgChan       chan interface{}
	AuthCache     map[string]BotAuth
}

func NewClientManager(config *common.Config) *ClientManager {
	return &ClientManager{
		ServerConfigs: config.Servers,
		AuthConfigs:   config.Auth,
		Clients:       make(map[string]*bot.Client),
		MsgChan:       make(chan interface{}),
		AuthCache:     make(map[string]BotAuth),
	}
}

// InitBotAuth 初始化机器人验证信息
func (c *ClientManager) InitBotAuth() {
	for name, auth := range c.AuthConfigs {
		if auth.Online {
			msAuth, err := GetMCcredentials("data/auth_"+name+".json", "88650e7e-efee-4857-b9a9-cf580a00ef43")
			if err != nil {
				common.Logger.Errorf("获取 %s 的验证信息失败: %s，将使用离线账户: %s", auth.Name, err, defaultName)
			}
			common.Logger.Infof("获取 %s 的验证信息成功: %s(%s)", name, msAuth.Name, msAuth.UUID)
			c.AuthCache[name] = msAuth
		} else {
			common.Logger.Infof("离线账户设置成功: %s", defaultName)
			c.AuthCache[name] = BotAuth{
				Name: defaultName,
			}
		}
	}
}

// JoinAllServers 加入所有服务器，出现错误按照配置的间隔尝试重连
func (c *ClientManager) JoinAllServers() {
	for _, server := range c.ServerConfigs {
		go c.JoinServer(server)
		time.Sleep(3 * time.Second)
	}

}

func (c *ClientManager) JoinServer(serverConfig common.ServerConfig) error {
	client := bot.NewClient()
	// 当缓存中没有验证信息时，使用默认离线账户
	auth, ok := c.AuthCache[serverConfig.Auth]
	if !ok {
		common.Logger.Warnf("未配置 %s 的验证信息，将使用离线账户 %s 加入 %s(%s)", serverConfig.Auth, defaultName, serverConfig.Name, serverConfig.Address)
		client.Auth.Name = defaultName
	} else {
		client.Auth.Name = auth.Name
		authConfig := c.AuthConfigs[serverConfig.Auth]
		// 如果是在线账户，将验证信息加入
		if authConfig.Online {
			client.Auth.UUID = auth.UUID
			client.Auth.AsTk = auth.AsTk
		} else {
		}
		common.Logger.Infof("将使用 %s 的验证信息 %s(%s) 加入 %s(%s)", serverConfig.Auth, auth.Name, auth.UUID, serverConfig.Name, serverConfig.Address)
	}
	basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{})
	err := client.JoinServer(serverConfig.Address)
	if err != nil {
		common.Logger.Errorf("使用 %s(%s) 加入服务器 %s 失败: %s", client.Auth.Name, client.Auth.UUID, serverConfig.Address, err)
		return err
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
	c.InitBotAuth()
	c.JoinAllServers()
}
