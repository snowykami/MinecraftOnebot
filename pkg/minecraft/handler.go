package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"strconv"
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
	pm, err := FormatPlayerMessage(msg, c.messageRegexps)
	if err != nil {
		c.GameLogInfof("%s: %s", common.Yellow("System"), msg)
		var event common.Event
		if strings.HasSuffix(msg.ClearString(), "joined the game") {
			Username := strings.Split(msg.ClearString(), " ")[0]
			event = common.Event{
				Type:       "notice",
				DetailType: "group_member_increase",
				UserID:     Username,
				GroupID:    c.Name,
				SubType:    "join",
				OperatorID: SystemUser,
				Data:       map[string]interface{}{},
			}
		} else if strings.HasSuffix(msg.ClearString(), "left the game") {
			Username := strings.Split(msg.ClearString(), " ")[0]
			event = common.Event{
				Type:       "notice",
				DetailType: "group_member_decrease",
				UserID:     Username,
				GroupID:    c.Name,
				OperatorID: SystemUser,
				SubType:    "leave",
				Data:       map[string]interface{}{},
			}
		} else {
			event = common.Event{
				Type:       "message",
				DetailType: "group",
				GroupID:    c.Name,
				UserID:     "System",
				Message:    msg.ClearString(),
			}
		}
		// 事件上行通道
		c.eventChan <- event
	} else {
		// 成功匹配为玩家消息
		// 推送
		var event common.Event
		if InArray(pm.Type, c.ServerConfig.PrivatePrefix) {
			event = common.Event{
				Type:       "message",
				DetailType: "private",
				UserID:     pm.Username,
				GroupID:    c.Name,
				Message:    pm.Message,
			}
		} else {
			event = common.Event{
				Type:       "message",
				DetailType: "group",
				UserID:     pm.Username,
				GroupID:    c.Name,
				UserTitle:  pm.Title,
				Message:    pm.Message,
			}
		}
		c.eventChan <- event
		// 日志输出
		if c.ServerConfig.IgnoreSelf && pm.Username == c.BotAuth.Name {
			return nil
		}
		msgText := ""
		if InArray(pm.Type, c.ServerConfig.PrivatePrefix) {
			msgText += common.Blue("私聊") + " "
		} else {
			msgText += common.Blue("公共") + " "
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

func (c *Connection) HandleConnectEvents(toManagerChan chan common.Event) {
	// 连接把事件传给管理器
	go func() {
		for {
			select {
			case event := <-toManagerChan:
				common.Logger.Infof("监听器接收到事件: %v\n", event)
			}
		}
	}()
	for {
		event := <-c.eventChan
		common.Logger.Infof("发送给Bot事件: %v", event.Message)
		toManagerChan <- event
	}
}

//func (m *Manager) HandleEvents() {
//	// 把事件传给Onebot
//	for {
//		event := <-m.ConnectionChan
//		common.Logger.Infof("发送给Bot事件: %v", event.Message)
//		m.EventChan <- event
//	}
//}

func (m *Manager) GetConnectionByIntIDStr(intID string) (*Connection, error) {
	for _, v := range m.Connections {
		if strconv.FormatInt(v.IntID, 10) == intID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("未找到服务器: %s", intID)
}

func (m *Manager) HandleEventsMux() {
	for {
		event := <-m.EventChan
		if event.Version == 11 {
			if event.Type == "message" {
				if event.SubType == "group" {
					// 获取serverName，并检测是否非nil
					connection, err := m.GetConnectionByIntIDStr(event.GroupID)
					if connection == nil {
						common.Logger.Warnln(err)
					} else {
						err := connection.chatHandler.SendMessage(event.Message)
						if err != nil {
							common.Logger.Warnln(err)
						}
					}
					// 传递事件
				} else if event.SubType == "private" {
					// 获取serverName，并检测是否非nil
					connection, err := m.GetConnectionByIntIDStr(event.GroupID)
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
				if event.SubType == "group" {
					// 获取serverName，并检测是否非nil
					connection, ok := m.Connections[event.GroupID]
					if !ok {
						common.Logger.Warnf("未找到服务器: %s", event.GroupID)
					} else {
						err := connection.chatHandler.SendMessage(event.Message)
						if err != nil {
							common.Logger.Warnln(err)
						}
					}
				} else if event.SubType == "private" {
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
