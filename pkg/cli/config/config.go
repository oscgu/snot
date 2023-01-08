package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

var Conf Config

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func Init(path string) {
	Conf.ParseOrDefault(path)
}

func (config *Config) ParseOrDefault(path string) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		createDefaultConf(path)
	} else {
		decoder := yaml.NewDecoder(f)
		if err = decoder.Decode(config); err != nil {
			handleErr(err)
		}

		decoder.Decode(config)
	}

	defer f.Close()
}

func createDefaultConf(path string) {
	conf := Config{
		User: User{
			Name:  "test-user",
			Group: "test",
		},
		Server: Server{
			Address: "",
			Port:    "",
			Active:  false,
		},
	}

	data, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println(err)
	}

	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		fmt.Println(err)
	}
}
