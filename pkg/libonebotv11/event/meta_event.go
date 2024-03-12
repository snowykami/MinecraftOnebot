package event

type MetaEvent struct {
	BaseEvent
	MetaEventType string `json:"meta_event_type"` // 元事件类型
	SubType       string `json:"sub_type"`        // 元事件子类型
}
