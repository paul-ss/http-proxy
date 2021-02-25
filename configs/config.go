package config

type Config struct {
	ProxyAddress string
}

var C = Config{
	ProxyAddress: "127.0.0.1:8000",
}
