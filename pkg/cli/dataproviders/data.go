package data

import (
	"github.com/oscgu/snot/pkg/cli/config"
	"github.com/oscgu/snot/pkg/cli/dataproviders/snotdb"
)

func GetProvider() config.DataProvider {
	if config.Conf.Server.Active {
		return nil
	}

	return snotdb.Snotdb
}
