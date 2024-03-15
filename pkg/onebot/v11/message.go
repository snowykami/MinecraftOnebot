package v11

import "encoding/json"

// 目前能在Minecraft中实现的消息段类型
const (
	MessageSegmentTypeText  = "text"
	MessageSegmentTypeAt    = "at"
	MessageSegmentTypeShare = "share"
	MessageSegmentTypeJson  = "json"
)

// MessageSegment 消息段
type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	Raw  string                 `json:"-"`
}

func NewMessageSegmentText(text string) *MessageSegment {
	return &MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
		Raw: text,
	}
}

// NewMessageSegmentAt at 消息段(原版实现)
func NewMessageSegmentAt(userID string) *MessageSegment {
	return &MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"user_id": userID,
		},
		Raw: "@" + userID,
	}
}

// NewMessageSegmentShare 分享消息段(可点击链接实现)
func NewMessageSegmentShare(url, title string) *MessageSegment {
	return &MessageSegment{
		Type: "share",
		Data: map[string]interface{}{
			"url":     url,
			"title":   title,
			"content": nil,
			"image":   nil,
		},
		Raw: url,
	}
}

// NewMessageSegmentJson json 消息段(tellraw实现)
func NewMessageSegmentJson(data map[string]interface{}) *MessageSegment {
	return &MessageSegment{
		Type: "json",
		Data: data,
		Raw: func() string {
			raw, err := json.Marshal(data)
			if err != nil {
				return ""
			}
			return string(raw)
		}(),
	}
}

func GetRawMessage(segments []*MessageSegment) string {
	var rawMessage string
	for _, segment := range segments {
		rawMessage += segment.Raw
	}
	return rawMessage
}
