package config

type Config struct {
	Addr  string
	Disks []string
}

func NewConfig(configPath string) *Config {
	return &Config{
		Addr: "0.0.0.0:9001",
		Disks: []string{
			"/tmp/disk1",
			"/tmp/disk2",
		},
	}
}
