package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/net"
	"sync/atomic"
	"time"
)

type RCONClient struct {
	Address   string
	Password  string
	conn      net.RCONClientConn
	closeChan chan int
	IsAlive   atomic.Bool
}

// RunRCON 运行RCON客户端
func (c *Connection) RunRCON() {
	common.Logger.Infof("正在连接RCON")
	if c.enableRCON {
		go c.RconClient.Connect()
		var reason string
		for {
			switch <-c.RconClient.closeChan {

			// 0: 连接失败，1: 连接断开，-1: 正常关闭
			case 0:
				{
					reason = fmt.Sprintf("%s(%s) RCON连接失败",
						c.ServerConfig.Address, c.Name)
				}
			case 1:
				{
					reason = fmt.Sprintf("%s(%s) RCON连接断开",
						c.ServerConfig.Address, c.Name)
				}
			case -1:
				{
					common.Logger.Infof("%s(%s) RCON连接正常关闭", c.ServerConfig.Address, c.Name)
					return
				}
			default:
				{
					common.Logger.Warnf("未知的RCON关闭原因")
				}
			}
			c.RconClient.IsAlive.Store(false)
			if c.ServerConfig.ReconnectInterval != 0 {
				common.Logger.Warnf("%s，将在%d秒后重连", reason, c.ServerConfig.ReconnectInterval)
			} else {
				common.Logger.Warnf("%s，配置文件未设置重连", reason)
				return
			}
			time.Sleep(time.Duration(c.ServerConfig.ReconnectInterval) * time.Second)
			go c.RconClient.Connect()
		}
	}
}

func (rc *RCONClient) Connect() {
	conn, err := net.DialRCON(rc.Address, rc.Password)
	if err != nil {
		rc.closeChan <- 0
	} else {
		rc.conn = conn
		rc.IsAlive.Store(true)
		common.Logger.Infof("%s RCON连接成功", rc.Address)
		rc.conn.Cmd("say RCON连接成功")
	}
}

func (rc *RCONClient) SendCommand(cmd string) error {
	err := rc.conn.Cmd(cmd)
	if err != nil {
		return err
	}
	return nil
}
