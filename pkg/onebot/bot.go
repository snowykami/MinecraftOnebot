package onebot

import (
	"MCOnebot/pkg/common"
	libobv11 "MCOnebot/pkg/libonebotv11"
	libobv12 "MCOnebot/pkg/libonebotv12"
	"MCOnebot/pkg/minecraft"
	"fmt"
	"strconv"
	"time"
)

const (
	Platform = "minecraft"
)

type Manager struct {
	V11Bot    *libobv11.OneBot
	V12Bot    *libobv12.OneBot
	EventChan chan common.Event
}

func CreateBots(config common.Config) (*libobv11.OneBot, *libobv12.OneBot, error) {
	//intSelfID, err := strconv.ParseInt(config.Onebot.Bot.SelfID, 10, 64)
	v11Used := false // 用于判断是否有使用到配置
	v12Used := false // 用于判断是否有使用到配置
	v11BotConfig := &libobv11.Config{
		Heartbeat: libobv11.ConfigHeartbeat{
			Enabled:  true,
			Interval: uint32(config.Onebot.Bot.HeartbeatInterval * 1000),
		},
		Comm: libobv11.ConfigComm{
			HTTP:        make([]libobv11.ConfigCommHTTP, 0),
			HTTPWebhook: make([]libobv11.ConfigCommHTTPWebhook, 0),
			WS:          make([]libobv11.ConfigCommWS, 0),
			WSReverse:   make([]libobv11.ConfigCommWSReverse, 0),
		},
	}

	v12BotConfig := &libobv12.Config{
		Heartbeat: libobv12.ConfigHeartbeat{
			Enabled:  true,
			Interval: uint32(config.Onebot.Bot.HeartbeatInterval * 1000),
		},
		Comm: libobv12.ConfigComm{
			HTTP:        make([]libobv12.ConfigCommHTTP, 0),
			HTTPWebhook: make([]libobv12.ConfigCommHTTPWebhook, 0),
			WS:          make([]libobv12.ConfigCommWS, 0),
			WSReverse:   make([]libobv12.ConfigCommWSReverse, 0),
		},
	}

	for _, reverseWSConf := range config.Onebot.ReverseWebSocket {

		if reverseWSConf.ProtocolVersion == 11 {
			v11BotConfig.Comm.WSReverse = append(v11BotConfig.Comm.WSReverse, libobv11.ConfigCommWSReverse{
				URL:               reverseWSConf.Address,
				AccessToken:       reverseWSConf.AccessToken,
				ReconnectInterval: uint32(reverseWSConf.ReconnectInterval * 1000),
			})
			v11Used = true
		} else {
			v12BotConfig.Comm.WSReverse = append(v12BotConfig.Comm.WSReverse, libobv12.ConfigCommWSReverse{
				URL:               reverseWSConf.Address,
				AccessToken:       reverseWSConf.AccessToken,
				ReconnectInterval: uint32(reverseWSConf.ReconnectInterval * 1000),
			})
			v12Used = true
		}
	}

	for _, wsConf := range config.Onebot.WebSocket {
		if wsConf.ProtocolVersion == 11 {
			v11BotConfig.Comm.WS = append(v11BotConfig.Comm.WS, libobv11.ConfigCommWS{
				Host:        wsConf.Host,
				Port:        uint16(wsConf.Port),
				AccessToken: wsConf.AccessToken,
			})
			v11Used = true
		} else {
			v12BotConfig.Comm.WS = append(v12BotConfig.Comm.WS, libobv12.ConfigCommWS{
				Host:        wsConf.Host,
				Port:        uint16(wsConf.Port),
				AccessToken: wsConf.AccessToken,
			})
			v12Used = true
		}
	}

	for _, httpWebhookConf := range config.Onebot.HTTPWebhook {
		if httpWebhookConf.ProtocolVersion == 11 {
			v11BotConfig.Comm.HTTPWebhook = append(v11BotConfig.Comm.HTTPWebhook, libobv11.ConfigCommHTTPWebhook{
				URL:         httpWebhookConf.Address,
				AccessToken: httpWebhookConf.AccessToken,
			})
			v11Used = true
		} else {
			v12BotConfig.Comm.HTTPWebhook = append(v12BotConfig.Comm.HTTPWebhook, libobv12.ConfigCommHTTPWebhook{
				URL:         httpWebhookConf.Address,
				AccessToken: httpWebhookConf.AccessToken,
			})
			v12Used = true
		}
	}

	for _, httpConf := range config.Onebot.HTTP {
		if httpConf.ProtocolVersion == 11 {
			v11BotConfig.Comm.HTTP = append(v11BotConfig.Comm.HTTP, libobv11.ConfigCommHTTP{
				Host:        httpConf.Host,
				Port:        uint16(httpConf.Port),
				AccessToken: httpConf.AccessToken,
			})
			v11Used = true
		} else {
			v12BotConfig.Comm.HTTP = append(v12BotConfig.Comm.HTTP, libobv12.ConfigCommHTTP{
				Host:        httpConf.Host,
				Port:        uint16(httpConf.Port),
				AccessToken: httpConf.AccessToken,
			})
			v12Used = true
		}
	}

	var v11Bot *libobv11.OneBot
	var v12Bot *libobv12.OneBot
	if v11Used {
		//intSelfID, err := strconv.ParseInt(config.Onebot.Bot.SelfID, 10, 64)
		intSelfID := minecraft.GenerateUserID(config.Onebot.Bot.SelfID)
		//if err != nil {
		//	return nil, nil, err
		//}
		self := &libobv11.Self{
			Platform: Platform,
			UserID:   intSelfID,
		}
		v11Bot = libobv11.NewOneBot(Platform, self, v11BotConfig)
	} else {
		v11Bot = nil
	}
	if v12Used {
		self := &libobv12.Self{
			Platform: Platform,
			UserID:   config.Onebot.Bot.SelfID,
		}
		v12Bot = libobv12.NewOneBot(Platform, self, v12BotConfig)
	} else {
		v12Bot = nil
	}
	return v11Bot, v12Bot, nil

}

func NewBotManager(config common.Config) (*Manager, error) {
	v11Bot, v12Bot, err := CreateBots(config)
	if err != nil {
		return nil, err
	}
	return &Manager{
		V11Bot: v11Bot,
		V12Bot: v12Bot,
	}, nil
}

func (m *Manager) Run() {
	//go m.OpenEventChan()
	go m.HandleActionMux()
	if m.V11Bot != nil {
		common.Logger.Info("OneBot v11 启动中...")
		go m.V11Bot.Run()
	}
	if m.V12Bot != nil {
		common.Logger.Info("OneBot v12 启动中...")
		go m.V12Bot.Run()
	}

}

func (m *Manager) OpenEventChan() {
	var event common.Event
	for {
		event = <-m.EventChan
		common.Logger.Infof("Bot收到事件: %v", event.Message)
		m.Push(event)
	}
}

func (m *Manager) Push(event common.Event) {
	if m.V11Bot != nil {
		switch event.Type {
		case "notice":
			obEvent := libobv11.MakeNoticeEvent(time.Now(), fmt.Sprintf("%v", event.DetailType))
			m.V11Bot.Push(&obEvent)
		case "message":
			// 构造消息
			messageID := int32(time.Now().UnixNano() / 1e6)
			message := libobv11.Message{
				libobv11.TextSegment(event.Message),
			}
			userID := minecraft.GenerateUserID(event.UserID)
			if event.DetailType == "private" {
				event := libobv11.MakePrivateMessageEvent(
					messageID, message, fmt.Sprintf("%v", event.Data["message"]), userID, event.UserID)
				m.V11Bot.Push(&event)
			} else if event.DetailType == "group" {
				GroupID := minecraft.GenerateUserID(event.GroupID)
				event := libobv11.MakeGroupMessageEvent(
					time.Now(), messageID, message, fmt.Sprintf("%v", event.Data["message"]), GroupID, userID, event.UserID, event.UserTitle)
				m.V11Bot.Push(&event)
			}
		}

	}
	if m.V12Bot != nil {
		switch event.Type {
		case "notice":
			if event.DetailType == "group_member_increase" {
				obEvent := libobv12.MakeGroupMemberIncreaseNoticeEvent(
					time.Now(), event.GroupID, event.UserID, event.OperatorID)
				obEvent.SubType = event.SubType
				m.V12Bot.Push(&obEvent)
			} else if event.DetailType == "group_member_decrease" {
				obEvent := libobv12.MakeGroupMemberDecreaseNoticeEvent(
					time.Now(), event.GroupID, event.UserID, event.OperatorID)
				obEvent.SubType = event.SubType
				m.V12Bot.Push(&obEvent)
			}
		case "message":
			messageID := strconv.FormatInt(time.Now().UnixNano(), 10)
			message := libobv12.Message{
				libobv12.TextSegment(event.Message),
			}
			if event.DetailType == "private" {
				event := libobv12.MakePrivateMessageEvent(
					time.Now(), messageID, message, event.Message, event.UserID)
				m.V12Bot.Push(&event)
			} else if event.DetailType == "group" {
				event := libobv12.MakeGroupMessageEvent(
					time.Now(), messageID, message, event.Message, event.GroupID, event.UserID)
				m.V12Bot.Push(&event)
			}
		}
	}
}
