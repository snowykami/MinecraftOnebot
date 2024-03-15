package minecraft

import (
	"MCOnebot/pkg/common"
	"errors"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/playerlist"
	"sync/atomic"
	"time"
)

var defaultName = "MCOnebot"

// Manager 客户端管理器
type Manager struct {
	Config      common.Config
	Connections map[string]*Connection // 客户端连接列表 name/group_id -> client
	EventChan   chan *Event
	AuthCache   map[string]*bot.Auth
}

type Connection struct {
	Name         string
	ServerConfig common.ServerConfig
	BotAuth      *bot.Auth
	Client       *bot.Client
	Handler      *GameHandler
	IsAlive      atomic.Bool

	player     *basic.Player
	playerList *playerlist.PlayerList
}

// NewClientManager 创建一个新的客户端管理器
func NewClientManager(config common.Config) *Manager {
	return &Manager{
		Config:      config,
		Connections: make(map[string]*Connection),
		EventChan:   make(chan *Event),
		AuthCache:   make(map[string]*bot.Auth),
	}
}

// InitBotAuth 初始化机器人验证信息
func (m *Manager) InitBotAuth() {
	for authName, authConfig := range m.Config.Auth {
		if authConfig.Online {
			// 配置的在线账户
			msAuth, err := GetMCcredentials("data/"+authName+".token", "88650e7e-efee-4857-b9a9-cf580a00ef43")
			if err != nil {
				common.Logger.Errorf("在线账户 %s 验证失败，将配置为离线账户 %s: %s", authName, defaultName, err)
			} else {
				common.Logger.Infof("在线账户验证成功 %s(%s) ，认证名称 %s", msAuth.Name, msAuth.UUID, authName)
			}
			m.AuthCache[authName] = &bot.Auth{
				Name: msAuth.Name,
				UUID: msAuth.UUID,
				AsTk: msAuth.AsTk,
			}
		} else {
			// 配置的离线账户
			common.Logger.Infof("离线账户设置成功 %s ，认证名称 %s", authConfig.Name, authName)

			m.AuthCache[authName] = &bot.Auth{
				Name: authConfig.Name,
			}
		}
	}
}

// Run 运行客户端管理器
func (m *Manager) Run() {
	// 等待账户验证完成后才可以启动游戏
	m.InitBotAuth()
	// 初始化连接内容
	for serverName, serverConfig := range m.Config.Servers {
		// 当缓存中没有验证信息时，使用默认离线账户
		botAuth, ok := m.AuthCache[serverConfig.Auth]
		if !ok {
			common.Logger.Errorf("未找到 %s 的验证信息，将使用默认离线账户 %s", serverConfig.Auth, defaultName)
			botAuth = &bot.Auth{
				Name: defaultName,
			}
		}
		connection := &Connection{
			Name:         serverName,
			ServerConfig: serverConfig,
			BotAuth:      botAuth,
		}
		connection.IsAlive.Store(true)
		m.Connections[serverName] = connection
	}
	for _, connection := range m.Connections {
		go connection.Run()
		time.Sleep(time.Duration(m.Config.Common.JoinInterval) * time.Second)
	}
	select {}
}

// Run 运行客户端,客户端一旦开始运行,自动重连由客户端自行管理
func (c *Connection) Run() {
	c.InitConnection()
	// 循环检测IsAlive状态
	go c.Join()
	for {
		if c.IsAlive.Load() {
			time.Sleep(5 * time.Second)
		} else {
			common.Logger.Printf("%s 与 %s 的连接断开，将在%d秒后重连", c.BotAuth.Name, c.ServerConfig.Address, c.ServerConfig.ReconnectInterval)
			time.Sleep(time.Duration(c.ServerConfig.ReconnectInterval) * time.Second)
			c.Join()
		}
	}
}

func (c *Connection) InitConnection() {
	c.Client = bot.NewClient()
	c.Client.Auth = *c.BotAuth
	c.playerList = playerlist.New(c.Client)
	c.player = basic.NewPlayer(c.Client, basic.DefaultSettings, basic.EventsListener{
		GameStart:    c.Handler.OnGameStart,
		Disconnect:   c.Handler.OnDisconnect,
		HealthChange: c.Handler.OnHealthChange,
		Death:        c.Handler.OnDeath,
		Teleported:   c.Handler.OnTeleported,
	})

}

func (c *Connection) Join() {
	err := c.Client.JoinServer(c.ServerConfig.Address)
	if err != nil {
		common.Logger.Errorf("使用 %s 加入服务器 %s 失败，将在%d秒后重连: %s",
			c.Client.Auth.Name, c.ServerConfig.Address, c.ServerConfig.ReconnectInterval, err,
		)
		c.IsAlive.Store(false)
	} else {
		common.Logger.Printf("%s 成功加入 %s(%s)", c.Client.Auth.Name, c.ServerConfig.Address, c.Name)
		c.IsAlive.Store(true)
	}
	var pErr bot.PacketHandlerError
	for {
		if err = c.Client.HandleGame(); err == nil {
			c.IsAlive.Store(false)
			return
		}
		if errors.As(err, &pErr) {
			common.Logger.Fatalln("处理单个数据包错误:", pErr)
		} else {
			common.Logger.Errorln("处理游戏错误:", err)
			c.IsAlive.Store(false)
			return
		}
	}
}
