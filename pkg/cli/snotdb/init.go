package snotdb

import (
	"fmt"

	cli "github.com/oscgu/snot/pkg/cli/note"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() {
	var err error
	Db, err = gorm.Open(sqlite.Open("snot.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	Db.AutoMigrate(&cli.Note{})
}
