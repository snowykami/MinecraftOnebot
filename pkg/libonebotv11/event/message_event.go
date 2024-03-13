package event

import (
	"strconv"
	"sync/atomic"
)

var (
	LatestMessageID uint64 = 0
)

type PrivateMessageEvent struct {
	BaseEvent
	MessageType string        `json:"message_type"` // 元事件类型
	SubType     string        `json:"sub_type"`     // 元事件子类型
	MessageID   string        `json:"message_id"`   // 消息 ID
	UserID      string        `json:"user_id"`      // 用户 ID
	Message     Message       `json:"message"`      // 消息内容
	RawMessage  string        `json:"raw_message"`  // 原始消息内容
	Font        int32         `json:"font"`         // 字体
	Sender      PrivateSender `json:"sender"`       // 发送者信息
}

type GroupMessageEvent struct {
	BaseEvent
	MessageType string        `json:"message_type"` // 元事件类型
	SubType     string        `json:"sub_type"`     // 元事件子类型
	MessageID   string        `json:"message_id"`   // 消息 ID
	GroupID     string        `json:"group_id"`     // 群号
	UserID      string        `json:"user_id"`      // 用户 ID
	Anonymous   Anonymous     `json:"anonymous"`    // 匿名信息
	Message     Message       `json:"message"`      // 消息内容
	RawMessage  string        `json:"raw_message"`  // 原始消息内容
	Font        int32         `json:"font"`         // 字体
	Sender      PrivateSender `json:"sender"`       // 发送者信息
}

// MakeGroupMessageEvent 快速构造一个群聊消息事件
func MakeGroupMessageEvent(message Message, groupID string, userID string) GroupMessageEvent {
	return GroupMessageEvent{
		BaseEvent:   MakeEvent(PostTypeMessage),
		MessageType: MessageTypeGroup,
		SubType:     GroupMessageEventSubTypeNormal,
		MessageID:   strconv.FormatUint(atomic.AddUint64(&LatestMessageID, 1), 10),
		GroupID:     groupID,
		UserID:      userID,
		Anonymous:   Anonymous{},
		Message:     message,
		RawMessage:  message.ExtractText(),
		Font:        0,               // 这逼东西我也没用过，懒得实现了，但是
		Sender:      PrivateSender{}, // 通过UserID获取
	}
}

// MakePrivateMessageEvent 快速构造一个私聊消息事件
func MakePrivateMessageEvent(message Message, userID string) PrivateMessageEvent {
	return PrivateMessageEvent{
		BaseEvent:   MakeEvent(PostTypeMessage),
		MessageType: MessageTypePrivate,
		SubType:     PrivateMessageEventSubTypeFriend,
		MessageID:   strconv.FormatUint(atomic.AddUint64(&LatestMessageID, 1), 10),
		UserID:      userID,
		Message:     message,
		RawMessage:  message.ExtractText(),
		Font:        0,               // 这逼东西我也没用过，懒得实现了，但是
		Sender:      PrivateSender{}, // 通过UserID获取
	}
}
