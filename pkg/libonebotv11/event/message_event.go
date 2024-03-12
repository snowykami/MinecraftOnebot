package event

import "strconv"

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

func MakeGroupMessageEvent(selfId int64, subType string, groupID string, userID string, anonymous Anonymous, message Message, font int32, sender PrivateSender) GroupMessageEvent {
	LatestMessageID++
	return GroupMessageEvent{
		BaseEvent:   MakeEvent(selfId, PostTypeMessage),
		MessageType: MessageTypeGroup,
		SubType:     GroupMessageEventSubTypeNormal,
		MessageID:   strconv.FormatUint(LatestMessageID, 10),
		GroupID:     groupID,
		UserID:      userID,
		Anonymous:   anonymous,
		Message:     message,
		RawMessage:  message.ExtractText(),
		Font:        font,
		Sender:      sender,
	}
}
