package main

import (
	cli "github.com/oscgu/snot/cmd/cli"
	"github.com/oscgu/snot/pkg/cli/config"
	db "github.com/oscgu/snot/pkg/cli/snotdb"
)

func main() {
	db.Init()
	config.Init()

	cli.Execute()
}
