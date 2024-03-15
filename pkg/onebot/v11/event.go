package v11

const (
	PostTypeMessage = "message"
	PostTypeNotice  = "notice"
	PostTypeRequest = "request"
	PostTypeMeta    = "meta_event"
)

const (
	MessageTypePrivate = "private"
	MessageTypeGroup   = "group"
)

const (
	SubTypeFriend    = "friend"
	SubTypeGroup     = "group"
	SubTypeOther     = "other"
	SubTypeNormal    = "normal"
	SubTypeAnonymous = "anonymous"
	SubTypeNotice    = "notice"
)

type BaseEvent struct {
	Time     int64  `json:"time"`
	SelfID   int64  `json:"self_id"`
	PostType string `json:"post_type"`
}

// MessageEvent 基础消息事件
type MessageEvent struct {
	*BaseEvent
	MessageType string            `json:"message_type"`
	SubType     string            `json:"sub_type"`
	MessageID   int32             `json:"message_id"`
	UserID      int64             `json:"user_id"`
	Message     []*MessageSegment `json:"message"`
	RawMessage  string            `json:"raw_message"`
	Font        int32             `json:"font"`
}

// NoticeEvent 通知事件
type NoticeEvent struct {
	*BaseEvent
	NoticeType string `json:"notice_type"`
}

// RequestEvent 请求事件
type RequestEvent struct {
	*BaseEvent
}

// MetaEvent 元事件
type MetaEvent struct {
	*BaseEvent
}

// PrivateMessageEvent 私聊消息事件
type PrivateMessageEvent struct {
	*MessageEvent
	Sender *PrivateSender `json:"sender"`
}

// GroupMessageEvent 群聊消息事件
type GroupMessageEvent struct {
	*MessageEvent
	GroupID   int64        `json:"group_id"`
	Anonymous interface{}  `json:"anonymous"`
	Sender    *GroupSender `json:"sender"`
}

type GroupAdminNoticeEvent struct {
	*NoticeEvent
	SubType string `json:"sub_type"`
	GroupID int64  `json:"group_id"`
	UserID  int64  `json:"user_id"`
}

type GroupDecreaseNoticeEvent struct {
	*NoticeEvent
	SubType    string `json:"sub_type"`
	GroupID    int64  `json:"group_id"`
	UserID     int64  `json:"user_id"`
	OperatorID int64  `json:"operator_id"`
}

type GroupIncreaseNoticeEvent struct {
	*NoticeEvent
	SubType    string `json:"sub_type"`
	GroupID    int64  `json:"group_id"`
	UserID     int64  `json:"user_id"`
	OperatorID int64  `json:"operator_id"`
}

type GroupBanNoticeEvent struct {
	*NoticeEvent
	SubType    string `json:"sub_type"`
	GroupID    int64  `json:"group_id"`
	UserID     int64  `json:"user_id"`
	OperatorID int64  `json:"operator_id"`
	Duration   int64  `json:"duration"`
}

// Sender 发送者 不单独用
type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`
}

// PrivateSender 私聊发送者
type PrivateSender struct {
	Sender
}

// GroupSender 群聊发送者
type GroupSender struct {
	Sender
	Card  string `json:"card"`
	Area  string `json:"area"`
	Level string `json:"level"`
	Role  string `json:"role"`
	Title string `json:"title"`
}

// Anonymous 匿名用户
type Anonymous struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// NewPrivateMessageEvent 创建一个私聊消息事件
func NewPrivateMessageEvent(segments []*MessageSegment, sender *PrivateSender) *PrivateMessageEvent {
	return &PrivateMessageEvent{
		MessageEvent: &MessageEvent{
			BaseEvent: &BaseEvent{
				Time:     0, // 自动填充
				PostType: PostTypeMessage,
				SelfID:   0, // 自动填充
			},
			MessageType: MessageTypePrivate,
			SubType:     SubTypeFriend,
			MessageID:   0, // 自动填充
			UserID:      sender.UserID,
			Message:     segments,
			RawMessage:  GetRawMessage(segments),
			Font:        10, // 在MC中无用
		},
		Sender: sender,
	}
}

// NewGroupMessageEvent 创建一个群聊消息事件
func NewGroupMessageEvent(segments []*MessageSegment, sender *GroupSender, groupID int64, anonymous interface{}) *GroupMessageEvent {
	return &GroupMessageEvent{
		MessageEvent: &MessageEvent{
			BaseEvent: &BaseEvent{
				Time:     0, // 自动填充
				PostType: PostTypeMessage,
				SelfID:   0, // 自动填充
			},
			MessageType: MessageTypeGroup,
			SubType:     SubTypeNormal,
			MessageID:   0, // 自动填充
			UserID:      sender.UserID,
			Message:     segments,
			RawMessage:  GetRawMessage(segments),
			Font:        10, // 在MC中无用
		},
		GroupID:   groupID,
		Anonymous: anonymous,
		Sender:    sender,
	}
}
