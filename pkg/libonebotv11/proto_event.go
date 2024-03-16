// OneBot Connect - 数据协议 - 事件
// https://12.onebot.dev/connect/data-protocol/event/

package libonebot

import (
	"errors"
	"time"
)

// EventTypeXxx 表示 OneBot 标准定义的事件类型.
const (
	EventTypeMessage = "message"    // 消息事件
	EventTypeNotice  = "notice"     // 通知事件
	EventTypeRequest = "request"    // 请求事件
	EventTypeMeta    = "meta_event" // 元事件
)

// Event 包含所有类型事件的共同字段.
type Event struct {
	// lock       sync.RWMutex
	SelfID   int64  `json:"self_id"`   // 事件 ID, 构造时自动生成
	Time     int64  `json:"time"`      // 事件发生时间 (Unix 时间戳), 单位: 秒
	PostType string `json:"post_type"` // 事件类型
}

func makeEvent(postType string) Event {
	return Event{
		Time:     time.Now().Unix(),
		PostType: postType,
		SelfID:   0, // 推送时构建
	}
}

// AnyEvent 是所有事件对象共同实现的接口.
type AnyEvent interface {
	Name() string
	tryFixUp(self *Self) error
}

// Name 返回事件名称.
func (e *Event) Name() string {
	// e.lock.RLock()
	// defer e.lock.RUnlock()
	return e.PostType
}

func (e *Event) tryFixUp(self *Self) error {
	// e.lock.Lock()
	// defer e.lock.Unlock()
	e.SelfID = self.UserID
	if e.Time == 0 {
		return errors.New("`time` 字段值无效")
	}
	if e.PostType != EventTypeMessage && e.PostType != EventTypeNotice && e.PostType != EventTypeRequest && e.PostType != EventTypeMeta {
		return errors.New("`type` 字段值无效")
	}
	return nil
}

// 四种事件基本类型

// MetaEvent 表示一个元事件.
type MetaEvent struct {
	Event
	MetaEventType string `json:"meta_event_type"` // 元事件类型
}

// MakeMetaEvent 构造一个元事件.
func MakeMetaEvent(metaEventType string) MetaEvent {
	return MetaEvent{
		Event:         makeEvent(EventTypeMeta),
		MetaEventType: metaEventType,
	}
}

// MessageEvent 表示一个消息事件.
type MessageEvent struct {
	Event
	MessageID  int32     `json:"message_id"`  // 消息 ID
	Message    []Segment `json:"message"`     // 消息内容
	RawMessage string    `json:"raw_message"` // 消息内容的替代表示, 可为空
	Font       int32     `json:"font"`        // 字体
}

// MakeMessageEvent 构造一个消息事件.
func MakeMessageEvent(detailType string, messageID int32, message []Segment, raw_message string) MessageEvent {
	return MessageEvent{
		Event:      makeEvent(EventTypeMessage),
		MessageID:  messageID,
		Message:    message,
		RawMessage: raw_message,
	}
}

// NoticeEvent 表示一个通知事件.
type NoticeEvent struct {
	Event
}

// MakeNoticeEvent 构造一个通知事件.
func MakeNoticeEvent(time time.Time, detailType string) NoticeEvent {
	return NoticeEvent{
		Event: makeEvent(EventTypeNotice),
	}
}

// RequestEvent 表示一个请求事件.
type RequestEvent struct {
	Event
}

// MakeRequestEvent 构造一个请求事件.
func MakeRequestEvent(time time.Time, detailType string) RequestEvent {
	return RequestEvent{
		Event: makeEvent(EventTypeRequest),
	}
}
