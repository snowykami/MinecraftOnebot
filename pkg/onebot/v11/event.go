package v11

const (
	PostTypeMessage = "message"
	PostTypeNotice  = "notice"
	PostTypeRequest = "request"
	PostTypeMeta    = "meta_event"
)

type BaseEvent struct {
	Time     int64  `json:"time"`
	SelfID   int64  `json:"self_id"`
	PostType string `json:"post_type"`
}

// MessageEvent 基础消息事件
type MessageEvent struct {
	BaseEvent
}

// NoticeEvent 通知事件
type NoticeEvent struct {
	BaseEvent
}

// RequestEvent 请求事件
type RequestEvent struct {
	BaseEvent
}

// MetaEvent 元事件
type MetaEvent struct {
	BaseEvent
}

// PrivateMessageEvent 私聊消息事件
type PrivateMessageEvent struct {
	MessageEvent
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   int32         `json:"message_id"`
	UserID      int64         `json:"user_id"`
	Message     string        `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Font        int32         `json:"font"`
	Sender      PrivateSender `json:"sender"`
}

// GroupMessageEvent 群聊消息事件
type GroupMessageEvent struct {
	MessageEvent
	MessageType string      `json:"message_type"`
	SubType     string      `json:"sub_type"`
	GroupID     int64       `json:"group_id"`
	UserID      int64       `json:"user_id"`
	Anonymous   Anonymous   `json:"anonymous"`
	Message     Message     `json:"message"`
	RawMessage  string      `json:"raw_message"`
	Font        int32       `json:"font"`
	Sender      GroupSender `json:"sender"`
}
