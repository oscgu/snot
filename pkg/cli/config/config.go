package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/oscgu/snot/internal/log"
	"gopkg.in/yaml.v3"
)

var Conf Config

func Init() {
	Conf.ParseOrDefault()
}

func GetConfDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	fullDirPath := filepath.Join(homeDir, ".snot")
	if err := os.MkdirAll(fullDirPath, 0700); err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(fullpath, data, 0644); err != nil {
		log.Fatal(err)
	}
}
