package event

import "time"

const (
	PostTypeMessage = "message"    // 消息事件
	PostTypeNotice  = "notice"     // 通知事件
	PostTypeRequest = "request"    // 请求事件
	PostTypeMeta    = "meta_event" // 元事件
)

const (
	MessageTypePrivate = "private" // 私聊消息
	MessageTypeGroup   = "group"   // 群消息
)

const (
	SexMale    = "male"
	SexFemale  = "female"
	SexUnknown = "unknown"
)

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

const (
	PrivateMessageEventSubTypeFriend = "friend" // 好友消息
	PrivateMessageEventSubTypeGroup  = "group"  // 群临时会话消息
	PrivateMessageEventSubTypeOther  = "other"  // 其他消息
)

const (
	GroupMessageEventSubTypeNormal    = "normal"    // 正常消息
	GroupMessageEventSubTypeAnonymous = "anonymous" // 匿名消息
	GroupMessageEventSubTypeNotice    = "notice"    // 群公告
)

type MessageSegment struct {
}

type Message struct {
	Segments []MessageSegment `json:"segments"` // 消息段列表
}

type PrivateSender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 发送者昵称
	Sex      string `json:"sex"`      // 发送者性别
	Age      int32  `json:"age"`      // 发送者年龄
}

type GroupSender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 发送者昵称
	Card     string `json:"card"`     // 发送者群名片／备注
	Sex      string `json:"sex"`      // 发送者性别
	Age      int32  `json:"age"`      // 发送者年龄
	Area     string `json:"area"`     // 发送者地区
	Level    string `json:"level"`    // 发送者成员等级
	Role     string `json:"role"`     // 发送者角色
	Title    string `json:"title"`    // 发送者专属头衔
}

type Anonymous struct {
	ID   string `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户 flag，在调用禁言 API 时需要传入
}

type BaseEvent struct {
	Time     int64  `json:"time"`      // 事件发生的时间戳
	SelfId   int64  `json:"self_id"`   // 事件发生的机器人 QQ 号
	PostType string `json:"post_type"` // 上报类型
}

func (m *Message) ExtractText() string {
	return ""
}

func MakeEvent(postType string) BaseEvent {
	return BaseEvent{
		Time:     time.Now().Unix(),
		SelfId:   0, // Push的时候会自动填充
		PostType: postType,
	}
}
