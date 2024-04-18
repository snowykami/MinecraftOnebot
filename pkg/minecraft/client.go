package minecraft

import (
	"MCOnebot/pkg/common"
	"errors"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"time"
)

type Client struct {
	Config common.ServerConfig
	Auth   bot.Auth

	botClient   *bot.Client
	player      *basic.Player
	chatHandler *msg.Manager
	playerList  *playerlist.PlayerList

	stopConnect chan string
}

func (c *Client) Start() {

	c.botClient = bot.NewClient()
	c.botClient.Auth = c.Auth

	c.playerList = playerlist.New(c.botClient)
	c.player = basic.NewPlayer(c.botClient, basic.DefaultSettings, basic.EventsListener{})
	c.chatHandler = msg.New(c.botClient, c.player, c.playerList, EventHandler)
	fmt.Println(c.chatHandler, c.playerList, c.player, c.botClient)

	go c.Connect()
	for {
		select {
		case <-c.stopConnect:
			if c.Config.ReconnectInterval > 0 {
				common.Logger.Warnf("与服务器连接 %s 断开，%d秒后重连", c.Config.Address, c.Config.ReconnectInterval)
				time.Sleep(time.Duration(c.Config.ReconnectInterval) * time.Second)
				go c.Connect()
			} else {
				common.Logger.Warnf("与服务器连接 %s 断开，不重连", c.Config.Address)
				return
			}
		}
	}
}

func (c *Client) Connect() {
	// TODO: 连接服务器
	err := c.botClient.JoinServer(c.Config.Address)
	if err != nil {
		common.Logger.Errorf("连接服务器 %s 失败: %v", c.Config.Address, err)
		c.stopConnect <- err.Error()
	}

	common.Logger.Infof("连接服务器 %s 成功", c.Config.Address)
	c.chatHandler.SendMessage("我进来了")

	var perr bot.PacketHandlerError
	for {
		if err = c.botClient.HandleGame(); err == nil {
			c.stopConnect <- err.Error()
		}
		if errors.As(err, &perr) {
			common.Logger.Print(perr)
		} else {
			c.stopConnect <- err.Error()
		}
	}
}
