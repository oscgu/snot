package snotdb

import (
	"fmt"
	"os"

	"github.com/oscgu/snot/pkg/cli/config"
	"github.com/oscgu/snot/pkg/cli/note"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Snotdb *DbProvider

func Init() {
	confDir := config.GetConfDir()
	db, err := gorm.Open(sqlite.Open(confDir+"/snot.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	Snotdb = &DbProvider{
		Db: db,
	}

	Snotdb.Db.AutoMigrate(&note.Note{})
}

type DbProvider struct {
	Db *gorm.DB
}

func (snotdb *DbProvider) GetTopics() []string {
	var topics []string
	snotdb.Db.Table("notes").Distinct("topic").Scan(&topics)

	return topics
}

func (snotdb *DbProvider) GetTitles(topic string) []string {
	var titles []string
	snotdb.Db.Table("notes").Where("topic = ?", topic).Select("title").Find(&titles)

	return titles
}

func (snotdb *DbProvider) GetNote(topic string, title string) (note.Note, error) {
	var n note.Note
	if err := snotdb.Db.Table("notes").
		Where("topic = ?", topic).
		Where("title = ?", title).
		Find(&n).Error; err != nil {
		return n, err
	}

	return n, nil
}

func (snotdb *DbProvider) SaveNote(note note.Note) {
	snotdb.Db.Table("notes").Create(note)
}
