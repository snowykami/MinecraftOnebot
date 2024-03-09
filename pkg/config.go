package pkg

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Name    string
	Address string
	Online  bool
}
