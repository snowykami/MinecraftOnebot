package common

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config 配置
type Config struct {
	Common  CommConfig            `yaml:"common"`  // 通用配置
	Servers []ServerConfig        `yaml:"servers"` // 服务器列表
	Auth    map[string]AuthConfig `yaml:"auth"`    // 验证信息
	Onebot  OnebotConfig          `yaml:"onebot"`  // Onebot 配置列表
}

// CommConfig 通用配置
type CommConfig struct {
	JoinInterval int `yaml:"join_interval"` // 加入间隔，单位秒，建议长一点
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name                 string `yaml:"name"` // 服务器名称，自定义值，建议为string类型的数字
	Address              string `yaml:"address"`
	ReconnectInterval    int    `yaml:"reconnect_interval"`     // 重连间隔，单位秒，建议长一点
	Auth                 string `yaml:"auth"`                   // 验证信息
	PlayerMessageHandler string `yaml:"player_message_handler"` // 玩家消息处理器，这是为了兼容bug
}

// AuthConfig 验证信息
type AuthConfig struct {
	Online       bool   `yaml:"online"`
	Name         string `yaml:"name"`
	JoinInterval int    `yaml:"join_interval"` // 加入间隔，单位秒，建议长一点
}

// OnebotConfig Onebot 配置
type OnebotConfig struct {
	ReverseWS []ReverseWSConfig `yaml:"reverse_ws"` // 反向 WebSocket 配置列表
	ForwardWS []ForwardWSConfig `yaml:"forward_ws"` // 正向 WebSocket 配置列表
	HTTPPost  []HTTPPostConfig  `yaml:"http_post"`  // HTTP POST 配置(反向 HTTP)
	HTTP      []HTTPConfig      `yaml:"http"`       // HTTP 配置(正向 HTTP)
	Bot       BotConfig         `yaml:"bot"`        // 机器人配置
}

// BotConfig 本地机器人配置
type BotConfig struct {
	SelfID            int64  `yaml:"self_id"`            // 机器人 QQ 号
	HeartbeatInterval int    `yaml:"heartbeat_interval"` // 心跳间隔, 单位为秒
	PlayerIDType      string `yaml:"player_id_type"`     // 玩家号传输类型
}

// ReverseWSConfig 反向 WebSocket 配置
type ReverseWSConfig struct {
	Address           string `yaml:"address"`            // 服务器地址 ws://example.com:8080/onebot/
	ReconnectInterval int    `yaml:"reconnect_interval"` // 重连间隔, 单位为秒
	AccessToken       string `yaml:"access_token"`       // AccessToken 用于验证连接的令牌
	ProtocolVersion   int    `yaml:"protocol_version"`   // 协议版本11/12
}

// ForwardWSConfig 正向 WebSocket 配置
type ForwardWSConfig struct {
	Host            string `yaml:"host"`             // 监听主机地址
	Port            int    `yaml:"port"`             // 绑定端口
	AccessToken     string `yaml:"access_token"`     // AccessToken 用于验证连接的令牌
	ProtocolVersion int    `yaml:"protocol_version"` // 协议版本11/12
}

// HTTPPostConfig HTTP POST 配置
type HTTPPostConfig struct {
	Address         string `yaml:"address"`          // 上报地址
	AccessToken     string `yaml:"access_token"`     // AccessToken 用于验证连接的令牌
	ProtocolVersion int    `yaml:"protocol_version"` // 协议版本11/12
}

// HTTPConfig HTTP 配置
type HTTPConfig struct {
	Host            string `yaml:"host"`             // 监听主机地址
	Port            int    `yaml:"port"`             // 绑定端口
	AccessToken     string `yaml:"access_token"`     // AccessToken 用于验证连接的令牌
	ProtocolVersion int    `yaml:"protocol_version"` // 协议版本11/12
}

// LoadConfig 从文件加载配置
func LoadConfig(fileName string, config *Config) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		Logger.Error("读取配置文件失败: ", err)
		return err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		Logger.Error("解析配置文件失败: ", err)
		return err
	}
	Logger.Println("配置文件加载成功")
	return nil
}
