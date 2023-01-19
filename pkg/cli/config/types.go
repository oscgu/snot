package config

type Config struct {
	User   User   `yaml:"user"`
	Server Server `yaml:"server"`
}

type User struct {
	Name  string `yaml:"name"`
	Group string `yaml:"group"`
}

type Server struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Active  bool   `yaml:"active"`
}
