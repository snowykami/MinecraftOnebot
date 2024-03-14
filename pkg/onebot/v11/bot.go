package v11

type Bot struct {
	SelfID      int64         `json:"self_id"`
	Connections []interface{} `json:"connections"`
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Run() int64 {
	return 0
}

func (b *Bot) Stop() int64 {
	return 0
}

func (b *Bot) PushEvent() int64 {
	return 0
}
