package common

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config 配置文件结构
type Config struct {
	Minecraft MinecraftConfig `yaml:"minecraft"`
	Onebot    OnebotConfig    `yaml:"onebot"`
	Redis     RedisConfig     `yaml:"redis"`
}

type MinecraftConfig struct {
	Servers []ServerConfig `yaml:"servers"`
	Auths   []AuthConfig   `yaml:"auths"`
}

type ServerConfig struct {
	Address           string `yaml:"address"`
	ID                int64  `yaml:"id"`
	Auth              string `yaml:"auth"`
	IgnoreSelf        bool   `yaml:"ignore_self"`
	ReconnectInterval int    `yaml:"reconnect_interval"`
}

type AuthConfig struct {
	Name   string `yaml:"name"`
	Online bool   `yaml:"online"`
	Player string `yaml:"player"` // 离线模式下为玩家名，正版模式留空
}

type OnebotConfig struct {
	impls           []map[string]any `yaml:"implementations"` // 动态类型列表，根据类型解析
	Implementations []interface{}
}

type OnebotImplementation interface {
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ImplConfig Implementations
type ImplConfig struct {
	Type         string `yaml:"type"`
	AccessTokens string `yaml:"access_tokens"`
}

type HttpImpl struct {
	ImplConfig
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type HttpPostImpl struct {
	ImplConfig
	Address string `yaml:"address"`
}

type ForwardWebSocketImpl struct {
	ImplConfig
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ReverseWebSocketImpl struct {
	ImplConfig
	Address           string `yaml:"address"`
	ReconnectInterval int    `yaml:"reconnect_interval"`
}

// LoadConfig 从文件config.yml中加载配置
func LoadConfig() (Config, error) {
	// 读取配置文件
	file, err := os.ReadFile("config.yml")
	if err != nil {
		Logger.Error("读取配置文件失败: ", err)
		return Config{}, err
	}
	// 解析配置文件
	config := new(Config)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		Logger.Error("解析配置文件失败: ", err)
		return Config{}, err
	}

	return *config, nil
}
