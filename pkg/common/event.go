package common

type Event struct {
	Type       string
	SubType    string
	DetailType string
	GroupID    int64
	UserID     int64
	Username   string
	GroupName  string
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
