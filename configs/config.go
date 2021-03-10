package config

type Config struct {
	ProxyAddress string
	CertsDir string
	MaxInMemoryCerts int
}

var C = Config{
	ProxyAddress: "127.0.0.1:8000",
	CertsDir: "certs/",
	MaxInMemoryCerts: 200,
}
