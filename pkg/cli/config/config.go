package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var Conf Config

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func Init() {
	Conf.ParseOrDefault()
}

func (config *Config) ParseOrDefault() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err)
	}

	fullDirPath := filepath.Join(homeDir, ".snot")
	if err := os.MkdirAll(fullDirPath, 0700); err != nil {
		handleErr(err)
	}

	fullConfPath := filepath.Join(fullDirPath, "config.yml")

	f, err := os.Open(fullConfPath)
	if os.IsNotExist(err) {
		createDefaultConf(fullConfPath)
	} else {
		decoder := yaml.NewDecoder(f)
		if err = decoder.Decode(config); err != nil {
			handleErr(err)
		}

		decoder.Decode(config)
	}

	defer f.Close()
}

func createDefaultConf(fullpath string) {
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

	if err = ioutil.WriteFile(fullpath, data, 0644); err != nil {
		fmt.Println(err)
	}
}
