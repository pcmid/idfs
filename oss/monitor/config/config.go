package config

type Config struct {
	Addr     string
	RingPath string
}

func NewConfig(configPath string) *Config {
	return &Config{
		Addr:     "0.0.0.0:9000",
		RingPath: "/tmp/ring/ring",
	}
}
