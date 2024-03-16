package common

type Event struct {
	Type       string
	SubType    string
	DetailType string
	GroupID    string
	UserID     string
	OperatorID string
	MessageID  string
	UserTitle  string
	Message    string
	Version    int
	Data       map[string]interface{}
}

type PrivateMessage struct {
	SenderName string
	Message    string
}

type PublicMessage struct {
	SenderName string
	Message    string
}

type SystemMessage struct {
	Message string
}
