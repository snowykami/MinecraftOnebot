package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/chat"
)

var EventHandler = msg.EventsHandler{
	SystemChat:        onSystemMsg,
	PlayerChatMessage: onPlayerMsg,
	DisguisedChat:     onDisguisedMsg,
}

func onSystemMsg(msg chat.Message, overlay bool) error {
	var sender, message string
	_, _ = fmt.Sscanf(msg.ClearString(), "<%s> %s", &sender, &message)

	common.Logger.Printf("%s: %s", sender, message)
	return nil
}

func onPlayerMsg(msg chat.Message, validated bool) error {
	common.Logger.Printf("Player: %s", msg)
	return nil
}

func onDisguisedMsg(msg chat.Message) error {
	common.Logger.Printf("Disguised: %v", msg)
	return nil
}
