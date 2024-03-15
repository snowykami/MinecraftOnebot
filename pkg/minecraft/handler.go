package minecraft

import (
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/chat"
)

type GameHandler struct {
	chatHandler *msg.Manager
	eventChan   chan *Event
}

func (h *GameHandler) HandleEVents() {
	// TODO
}

func (h *GameHandler) OnGameStart() error {
	// TODO
	return nil
}

func (h *GameHandler) OnDisconnect(reason chat.Message) error {
	// TODO
	return nil
}

func (h *GameHandler) OnHealthChange(health float32, foodLevel int32, foodSaturation float32) error {
	// TODO
	return nil
}

func (h *GameHandler) OnDeath() error {
	// TODO
	return nil
}

func (h *GameHandler) OnTeleported(x, y, z float64, yaw, pitch float32, flags byte, teleportID int32) error {
	// TODO
	return nil
}

// MessageHandler

func (h *GameHandler) OnPlayerChatMessage(msg chat.Message, overlay bool) error {
	// TODO
	return nil
}

func (h *GameHandler) OnSystemChat(msg chat.Message, validated bool) error {
	// TODO
	return nil
}

func (h *GameHandler) OnDIsguisedChat(msg chat.Message) error {
	// TODO
	return nil
}
