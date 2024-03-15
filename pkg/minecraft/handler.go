package minecraft

import (
	"MCOnebot/pkg/common"
	"github.com/Tnze/go-mc/chat"
)

func (c *Connection) HandleEVents() {
	// TODO
}

func (c *Connection) OnGameStart() error {
	c.GameLogInfof("GameStart")
	// TODO
	return nil
}

func (c *Connection) OnDisconnect(reason chat.Message) error {
	c.GameLogWarnf("Disconnect: %s", reason)
	// TODO
	return nil
}

func (c *Connection) OnHealthChange(health float32, foodLevel int32, foodSaturation float32) error {
	c.GameLogInfof("HealthChange: %f, FoodLevel: %d, FoodSaturation: %f", health, foodLevel, foodSaturation)
	// TODO
	return nil
}

func (c *Connection) OnDeath() error {
	c.GameLogInfof("Death")
	// TODO
	return nil
}

func (c *Connection) OnTeleported(x, y, z float64, yaw, pitch float32, flags byte, teleportID int32) error {
	c.GameLogInfof("Teleported: %f, %f, %f, %f, %f, %d, %d", x, y, z, yaw, pitch, flags, teleportID)
	// TODO
	return nil
}

// MessageHandler

func (c *Connection) OnPlayerChatMessage(msg chat.Message, overlay bool) error {
	c.GameLogInfof("%s: %s", common.Blue("Player"), msg)
	// TODO
	return nil
}

func (c *Connection) OnSystemChat(msg chat.Message, validated bool) error {
	pm, err := FormatPlayerMessage(msg, c.messageRegexps)
	if err != nil {
		c.GameLogInfof("%s: %s", common.Yellow("System"), msg)
	} else {
		// 成功匹配为玩家消息
		if pm.Title != "" {
			c.GameLogInfof("[%s]%s: %s", common.Blue(pm.Title), common.Cyan(pm.Username), pm.Message)
			return nil
		} else {
			c.GameLogInfof("%s: %s", common.Cyan(pm.Username), pm.Message)
		}
	}
	// TODO
	return nil
}

func (c *Connection) OnDisguisedChat(msg chat.Message) error {
	c.GameLogInfof("Disguised: %s", msg)
	// TODO
	return nil
}
