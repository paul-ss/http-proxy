package config


type Config struct {
	ProxyAddress string
	ApiAddress string
	CertsDir string
	MaxInMemoryCerts int
	Db DbConf
}

type DbConf struct {
	User string
	Password string
	Host string
	Port string
	DbName string
}

var C = Config{
	ProxyAddress: "0.0.0.0:8000",
	ApiAddress: "0.0.0.0:8080",
	CertsDir: "certs/",
	MaxInMemoryCerts: 200,
	Db: DbConf{
		Host: "localhost",
		Port: "5432",
		DbName: "proxy_db",
		User: "proxy",
		Password: "jw8s0F4",
	},
}
