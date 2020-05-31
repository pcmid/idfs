package common

type Config struct {
	Addr    string
	Backend struct {
		Url string
	}
}

func NewConfig(configPath string) *Config {
	return &Config{
		Addr:    "127.0.0.1:8000",
		Backend: struct{ Url string }{Url: "http://127.0.0.1:9000/oss/"},
	}
}
