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

func GetConfDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err)
	}

	fullDirPath := filepath.Join(homeDir, ".snot")
	if err := os.MkdirAll(fullDirPath, 0700); err != nil {
		handleErr(err)
	}

	return fullDirPath
}

func (config *Config) ParseOrDefault() {
	confDir := GetConfDir()
	fullConfPath := filepath.Join(confDir, "config.yml")

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
		Editor: "default",
		Server: Server{
			Address: "",
			Port:    "",
			Active:  false,
		},
	}

	data, err := yaml.Marshal(conf)
	if err != nil {
		handleErr(err)
	}

	if err = ioutil.WriteFile(fullpath, data, 0644); err != nil {
		handleErr(err)
	}
}
