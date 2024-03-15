package common

import (
	"MCOnebot/pkg/common"
	v11 "MCOnebot/pkg/onebot/v11"
)

type WebsocketReverseComm struct {
	Bot    *v11.Bot
	Config *common.WebSocketConfig
}
