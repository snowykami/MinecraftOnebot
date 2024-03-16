package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"strings"
)

const SystemUser = "System"

func (c *Connection) OnGameStart() error {
	c.GameLogInfof("GameStart")
	// TODO
	return nil
}

func (c *Connection) OnDisconnect(reason chat.Message) error {
	c.GameLogWarnf("Disconnect: %s", reason)
	// TODO
	return nil
}

func (c *Connection) OnHealthChange(health float32, foodLevel int32, foodSaturation float32) error {
	c.GameLogInfof("HealthChange: %f, FoodLevel: %d, FoodSaturation: %f", health, foodLevel, foodSaturation)
	// TODO
	return nil
}

func (c *Connection) OnDeath() error {
	c.GameLogInfof("Death")
	// TODO
	return nil
}

func (c *Connection) OnTeleported(x, y, z float64, yaw, pitch float32, flags byte, teleportID int32) error {
	c.GameLogInfof("Teleported: %f, %f, %f, %f, %f, %d, %d", x, y, z, yaw, pitch, flags, teleportID)
	// TODO
	return nil
}

// MessageHandler

func (c *Connection) OnPlayerChatMessage(msg chat.Message, overlay bool) error {
	c.GameLogInfof("%s: %s", common.Blue("Player"), msg)
	// TODO
	return nil
}

func (c *Connection) OnSystemChat(msg chat.Message, validated bool) error {
	pm, err := FormatPlayerMessage(msg, c.messageRegexps, c.playerList)
	if err != nil {
		if c.ServerConfig.IgnoreSelf && msg.ClearString() == c.lastSentMessage {
			return nil
		}
		c.GameLogInfof("%s: %s", common.Yellow("System"), msg)
		var event common.Event
		if strings.HasSuffix(msg.ClearString(), "joined the game") {
			Username := strings.Split(msg.ClearString(), " ")[0]
			event = common.Event{
				Type:       "notice",
				DetailType: "group_member_increase",
				Username:   Username,
				UserID:     GenerateIntID(Username),
				GroupName:  c.Name,
				GroupID:    c.ID,
				SubType:    "join",
				OperatorID: SystemUser,
				Data:       map[string]interface{}{},
			}
		} else if strings.HasSuffix(msg.ClearString(), "left the game") {
			Username := strings.Split(msg.ClearString(), " ")[0]
			event = common.Event{
				Type:       "notice",
				DetailType: "group_member_decrease",
				Username:   Username,
				UserID:     GenerateIntID(Username),
				GroupName:  c.Name,
				GroupID:    c.ID,
				OperatorID: SystemUser,
				SubType:    "leave",
				Data:       map[string]interface{}{},
			}
		} else {
			event = common.Event{
				Type:       "message",
				DetailType: "group",
				GroupName:  c.Name,
				GroupID:    c.ID,
				Username:   "System",
				UserID:     GenerateIntID("System"),
				Message:    msg.ClearString(),
			}
		}
		// 事件上行通道
		c.eventChan <- event
	} else {
		// 成功匹配为玩家消息
		// 推送
		if c.ServerConfig.IgnoreSelf && RemoveANSI(pm.Username) == c.BotAuth.Name {
			return nil
		}
		var event common.Event
		if InArray(pm.Type, c.ServerConfig.PrivatePrefix) {
			event = common.Event{
				Type:       "message",
				DetailType: "private",
				Username:   RemoveANSI(pm.Username),
				UserID:     GenerateIntID(RemoveANSI(pm.Username)),
				GroupName:  c.Name,
				GroupID:    c.ID,
				Message:    RemoveANSI(pm.Message),
			}
		} else {
			event = common.Event{
				Type:       "message",
				DetailType: "group",
				Username:   RemoveANSI(pm.Username),
				UserID:     GenerateIntID(RemoveANSI(pm.Username)),
				GroupName:  c.Name,
				GroupID:    c.ID,
				UserTitle:  RemoveANSI(pm.Title),
				Message:    RemoveANSI(pm.Message),
			}
		}
		c.eventChan <- event
		// 日志输出
		msgText := ""
		if InArray(pm.Type, c.ServerConfig.PrivatePrefix) {
			msgText += common.Blue("私聊") + " "
		}
		if pm.Title != "" {
			msgText += fmt.Sprintf("[%s]", common.Cyan(pm.Title)) + " "
		}
		msgText += fmt.Sprintf("%s: %s", common.Yellow(pm.Username), pm.Message)
		c.GameLogInfof(msgText)
	}
	// TODO
	return nil
}

func (c *Connection) OnDisguisedChat(msg chat.Message) error {
	c.GameLogInfof("Disguised: %s", msg)
	// TODO
	return nil
}

// SendMessage 自动处理消息，用go执行
func (c *Connection) SendMessage(message string) error {
	for _, v := range c.illegalChar {
		message = strings.ReplaceAll(message, v, "")
	}
	c.lastSentMessage = message
	// 多行消息分开发送，先判断有无启用RCON，有则使用RCON发送tellraw
	if c.enableRCON && c.RconClient.IsAlive.Load() && c.RconClient.conn != nil {
		go func() {
			//err := c.RconClient.SendCommand(fmt.Sprintf("tellraw @a {\"text\":\"%s\"}", message))
			//if err != nil {
			//	return
			//}
			for _, v := range strings.Split(message, "\n") {
				common.Logger.Infof("RCON: %s", v)
				err := c.RconClient.conn.Cmd(fmt.Sprintf("say %s", v))
				if err != nil {
				}
			}
		}()
		return nil
	}
	for _, v := range strings.Split(message, "\n") {
		err := c.chatHandler.SendMessage(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection) HandleConnectEvents(toManagerChan chan common.Event) {
	for {
		select {
		case event := <-c.eventChan:
			toManagerChan <- event
		}
	}
}

func (m *Manager) GetConnectionByIntID(intID int64) (*Connection, error) {
	for _, v := range m.Connections {
		if v.ID == intID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("未找到服务器: %d", intID)
}

func (m *Manager) GetConnectionByName(name string) (*Connection, error) {
	for _, v := range m.Connections {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, fmt.Errorf("未找到服务器: %s", name)

}

func (m *Manager) HandleEventsMux() {
	for {
		event := <-m.SendChan
		if event.Version == 11 {
			if event.Type == "message" {

				if event.DetailType == "group" {
					// 获取serverName，并检测是否非nil
					connection, err := m.GetConnectionByIntID(event.GroupID)
					if connection == nil {
						common.Logger.Warnln(err)
					} else {
						err := connection.SendMessage(event.Message)
						if err != nil {
							common.Logger.Warnln(err)
						}
					}
					// 传递事件
				} else if event.DetailType == "private" {
					// 获取serverName，并检测是否非nil
					connection, err := m.GetConnectionByIntID(event.GroupID)
					if connection == nil {
						common.Logger.Warnln(err)
					} else {
						//err := connection.chatHandler.SendPrivateMessage(event.UserID, event.Message)
						//if err != nil {
						//	common.Logger.Warnln(err)
						//}
					}
				}
			}
		} else if event.Version == 12 {
			// 12无需处理ID，原始ID即为serverName
			if event.Type == "message" {

				if event.DetailType == "group" {
					// 获取serverName，并检测是否非nil
					connection, err := m.GetConnectionByName(event.GroupName)
					if err != nil {
						common.Logger.Warnf("未找到服务器: %s", event.GroupName)
					} else {
						err := connection.SendMessage(event.Message)
						if err != nil {
							common.Logger.Warnln(err)
						}
					}
				} else if event.DetailType == "private" {
					//// 获取serverName，并检测是否非nil
					//connection, ok := m.Connections[event.GroupID]
					//if !ok {
					//	common.Logger.Warnf("未找到服务器: %s", event.GroupID)
					//} else {
					//	//err := connection.chatHandler.SendPrivateMessage(event.UserID, event.Message)
					//	//if err != nil {
					//	//	common.Logger.Warnln(err)
					//	//}
					//}
				}
			}
		}
	}
}
