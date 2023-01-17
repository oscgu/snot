package snotdb

import (
	"fmt"
	"os"

	c "github.com/oscgu/snot/pkg/cli/config"
	cli "github.com/oscgu/snot/pkg/cli/note"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() {
	var err error
	confDir := c.GetConfDir()
	Db, err = gorm.Open(sqlite.Open(confDir+"/snot.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	Db.AutoMigrate(&cli.Note{})
}
