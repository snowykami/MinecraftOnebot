package v11

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
