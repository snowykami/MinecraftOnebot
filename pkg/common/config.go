package common

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config 配置
type Config struct {
	Common  CommConfig              `yaml:"common"`  // 通用配置
	Servers map[string]ServerConfig `yaml:"servers"` // 服务器列表
	Auth    map[string]AuthConfig   `yaml:"auth"`    // 验证信息
	Onebot  OnebotConfig            `yaml:"onebot"`  // Onebot 配置列表
}

// CommConfig 通用配置
type CommConfig struct {
	JoinInterval int `yaml:"join_interval"` // 加入间隔，单位秒，建议长一点
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address           string     `yaml:"address"`
	ReconnectInterval int        `yaml:"reconnect_interval"` // 重连间隔，单位秒，建议长一点
	Auth              string     `yaml:"auth"`               // 验证信息
	MessageTemplates  []string   `yaml:"message_templates"`  // 玩家消息处理器，这是为了兼容bug
	PrivatePrefix     []string   `yaml:"private_prefix"`     // 私聊前缀
	IgnoreSelf        bool       `yaml:"ignore_self"`        // 忽略自己的消息
	RCON              RCONConfig `yaml:"rcon"`               // RCON 配置
}

// RCONConfig RCON 配置
type RCONConfig struct {
	Address  string `yaml:"address"`  // RCON服务器地址
	Password string `yaml:"password"` // RCON 密码
}

// AuthConfig 验证信息
type AuthConfig struct {
	Online       bool   `yaml:"online"`
	Name         string `yaml:"name"`
	JoinInterval int    `yaml:"join_interval"` // 加入间隔，单位秒，建议长一点
}

// OnebotConfig Onebot 配置
type OnebotConfig struct {
	ReverseWebSocket []ReverseWebSocketConfig `yaml:"reverse_websocket"` // 反向 WebSocket 配置列表
	WebSocket        []WebSocketConfig        `yaml:"websocket"`         // 正向 WebSocket 配置列表
	HTTPWebhook      []HTTPWebhookConfig      `yaml:"http_webhook"`      // HTTP POST 配置(反向 HTTP)
	HTTP             []HTTPConfig             `yaml:"http"`              // HTTP 配置(正向 HTTP)
	Bot              BotConfig                `yaml:"bot"`               // 机器人配置
}

// BotConfig 本地机器人配置
type BotConfig struct {
	SelfID            string `yaml:"self_id"`            // 机器人 QQ 号
	HeartbeatInterval int    `yaml:"heartbeat_interval"` // 心跳间隔, 单位为秒
	PlayerIDType      string `yaml:"player_id_type"`     // 玩家号传输类型
}

// ReverseWebSocketConfig 反向 WebSocket 配置
type ReverseWebSocketConfig struct {
	Address           string `yaml:"address"`            // 服务器地址 ws://example.com:8080/onebot/
	ReconnectInterval int    `yaml:"reconnect_interval"` // 重连间隔, 单位为秒
	AccessToken       string `yaml:"access_token"`       // AccessToken 用于验证连接的令牌
	ProtocolVersion   int    `yaml:"protocol_version"`   // 协议版本11/12
}

// WebSocketConfig 正向 WebSocket 配置
type WebSocketConfig struct {
	Host            string `yaml:"host"`             // 监听主机地址
	Port            int    `yaml:"port"`             // 绑定端口
	AccessToken     string `yaml:"access_token"`     // AccessToken 用于验证连接的令牌
	ProtocolVersion int    `yaml:"protocol_version"` // 协议版本11/12
}

// HTTPWebhookConfig HTTP POST 配置
type HTTPWebhookConfig struct {
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
	_, err := os.Stat(fileName)
	if err != nil {
		Logger.Error("配置文件 config.yml 不存在，请修改 config.example.yml 后自行保存为 config.yml 后重启", err)
		return err
	}
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
