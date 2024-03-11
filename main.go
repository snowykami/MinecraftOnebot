package main

import (
	libob "MCOnebot/pkg/libonebotv11"
	"time"
)

var (
	ob *libob.OneBot
)

func main() {
	self := &libob.Self{
		Platform: "minecraft",
		UserID:   "Steve",
	}
	config := &libob.Config{
		Heartbeat: libob.ConfigHeartbeat{
			Enabled:  true,
			Interval: 20000,
		},
		Comm: libob.ConfigComm{
			WSReverse: []libob.ConfigCommWSReverse{
				{
					URL:               "ws://127.0.0.1:20216/onebot/v11/ws",
					ReconnectInterval: 3000,
				},
			},
		},
	}

	ob = libob.NewOneBot("minecraft-onebot", self, config)
	go ob.Run()
	time.Sleep(1 * time.Second)
	msg := libob.Message{
		libob.TextSegment("Hello, World!"),
	}
	groupEvent := libob.MakeGroupMessageEvent(time.Now(), "1", msg, "alt_message", "lys", "Steve")
	go ob.Push(&groupEvent)
	time.Sleep(5 * time.Second)
}
