package config

type ServerConfig struct {
	App     app     `yaml:"app"`
	Elastic elastic `yaml:"elastic"`
}

type app struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type elastic struct {
	Url   string `yaml:"url"`
	Sniff bool   `yaml:"sniff"`
}
