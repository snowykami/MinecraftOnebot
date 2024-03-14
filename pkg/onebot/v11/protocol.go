package v11

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

package v11

import (
"fmt"
"log"
"github.com/gorilla/websocket"
)

// RunWebSocketClient 创建WebSocket客户端，反向WS
func RunWebSocketClient(config *common.ReverseWSConfig) (*websocket.Conn, error) {
	headers := http.Header{}
	if config.AccessToken != "" {
		headers.Add("Authorization", fmt.Sprintf("Bearer %s", config.AccessToken))
	}
	conn, _, err := websocket.DefaultDialer.Dial(config.Address, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conn, nil
}