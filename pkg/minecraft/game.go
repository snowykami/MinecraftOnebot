package minecraft

import (
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
)

func JoinServer() {
	client := bot.NewClient()
	client.Auth = basic.Auth("email", "password")
}
