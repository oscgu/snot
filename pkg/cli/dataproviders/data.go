package data

import (
	"github.com/oscgu/snot/pkg/cli/config"
	"github.com/oscgu/snot/pkg/cli/dataproviders/snotdb"
	"github.com/oscgu/snot/pkg/cli/note"
)

type DataProvider interface {
	GetTopics() []string
	GetTitles(topic string) []string
	GetNote(topic string, title string) (note.Note, error)
}

func GetProvider() DataProvider {
	if config.Conf.Server.Active {
		return nil
	}

	return snotdb.Snotdb
}
