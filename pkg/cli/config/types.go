package config

import "github.com/oscgu/snot/pkg/cli/note"

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

type DataProvider interface {
	GetTopics() []string
	GetTitles(topic string) []string
	GetNote(topic string, title string) (note.Note, error)
}
