package minecraft

import (
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
)

var (
	client *bot.Client
	player *basic.Player
)

func JoinServer() {
	client = bot.NewClient()
}
