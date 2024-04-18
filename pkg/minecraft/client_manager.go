package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/bot"
)

type ClientManager struct {
	Config  common.MinecraftConfig
	Clients map[int64]*Client
}

func NewClientManager(config common.MinecraftConfig) *ClientManager {
	return &ClientManager{
		Config:  config,
		Clients: make(map[int64]*Client),
	}
}

func (cm *ClientManager) Start() {
	// 获取验证信息
	auths := make(map[string]bot.Auth)
	auths[""] = bot.Auth{} // 空验证信息

	// 鉴权验证
	for _, a := range cm.Config.Auths {
		var onlineType = "在线"
		if a.Online {
			var botAuth bot.Auth
			botAuth, err := GetMCcredentials(fmt.Sprintf("data/%s.token", a.Name), "88650e7e-efee-4857-b9a9-cf580a00ef43")
			if err != nil {
				onlineType = "离线"
				common.Logger.Errorf("Failed to load auth %s: %s", a.Name, err)
				botAuth = bot.Auth{
					Name: "Liteyuki",
				}
			}
			auths[a.Name] = botAuth
		} else {
			onlineType = "离线"
			auths[a.Name] = bot.Auth{
				Name: a.Player,
			}
		}
		common.Logger.Infof("成功加载验证信息: %s(%s %s)", a.Player, a.Name, onlineType)
	}

	// 初始化和启动客户端
	for _, c := range cm.Config.Servers {
		// 初始化客户端
		if c.ID == 0 {
			common.Logger.Error("服务器ID不能为0")
			continue
		}
		if c.Address == "" {
			common.Logger.Error("服务器地址不能为空")
			continue
		}
		cm.Clients[c.ID] = &Client{
			Config:      c,
			Auth:        auths[c.Auth],
			stopConnect: make(chan string),
		}
		// 启动客户端
		go cm.Clients[c.ID].Start()
	}
}
