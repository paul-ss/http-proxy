package config

type Config struct {
	ProxyAddress string
	CertsDir string
	MaxInMemoryCerts int
}

var C = Config{
	ProxyAddress: "0.0.0.0:8000",
	CertsDir: "certs/",
	MaxInMemoryCerts: 200,
}
