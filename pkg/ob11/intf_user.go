// 接口定义 - 单用户接口

package libonebot

import "time"

// 用户消息事件
// https://12.onebot.dev/interface/user/message-events/

// PrivateMessageEvent 表示一个私聊消息事件.
type PrivateMessageEvent struct {
	MessageEvent
	UserID      int64                  `json:"user_id"`      // 用户 ID
	MessageType string                 `json:"message_type"` // 消息类型
	SubType     string                 `json:"sub_type"`     // 消息子类型
	Sender      map[string]interface{} `json:"sender"`       // 发送者信息
}

// MakePrivateMessageEvent 构造一个私聊消息事件.
func MakePrivateMessageEvent(messageID int32, message []Segment, altMessage string, userID int64, userName string) PrivateMessageEvent {
	return PrivateMessageEvent{
		MessageEvent: MakeMessageEvent("private", messageID, message, altMessage),
		UserID:       userID,
		MessageType:  "private",
		SubType:      "friend",
		Sender: map[string]interface{}{
			"user_id":  userID,
			"nickname": userName,
			"age":      0,
			"sex":      "unknown",
		},
	}
}

// 用户通知事件
// https://12.onebot.dev/interface/user/notice-events/

// FriendIncreaseNoticeEvent 表示一个好友增加通知事件.
type FriendIncreaseNoticeEvent struct {
	NoticeEvent
	UserID string `json:"user_id"` // 用户 ID
}

// MakeFriendIncreaseNoticeEvent 构造一个好友增加通知事件.
func MakeFriendIncreaseNoticeEvent(time time.Time, userID string) FriendIncreaseNoticeEvent {
	return FriendIncreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "friend_increase"),
		UserID:      userID,
	}
}

// FriendDecreaseNoticeEvent 表示一个好友减少通知事件.
type FriendDecreaseNoticeEvent struct {
	NoticeEvent
	UserID string `json:"user_id"` // 用户 ID
}

// MakeFriendDecreaseNoticeEvent 构造一个好友减少通知事件.
func MakeFriendDecreaseNoticeEvent(time time.Time, userID string) FriendDecreaseNoticeEvent {
	return FriendDecreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "friend_decrease"),
		UserID:      userID,
	}
}

// PrivateMessageDeleteNoticeEvent 表示一个私聊消息删除通知事件.
type PrivateMessageDeleteNoticeEvent struct {
	NoticeEvent
	MessageID string `json:"message_id"` // 消息 ID
	UserID    string `json:"user_id"`    // 消息发送者 ID
}

// MakePrivateMessageDeleteNoticeEvent 构造一个私聊消息删除通知事件.
func MakePrivateMessageDeleteNoticeEvent(time time.Time, messageID string, userID string) PrivateMessageDeleteNoticeEvent {
	return PrivateMessageDeleteNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "private_message_delete"),
		MessageID:   messageID,
		UserID:      userID,
	}
}

// 用户动作
// https://12.onebot.dev/interface/user/actions/

const (
	ActionGetSelfInfo   = "get_self_info"   // 获取机器人自身信息
	ActionGetUserInfo   = "get_user_info"   // 获取用户信息
	ActionGetFriendList = "get_friend_list" // 获取好友列表
)
