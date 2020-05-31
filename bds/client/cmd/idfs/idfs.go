package main

import (
	log "github.com/sirupsen/logrus"
	"idfs/bds/common"
	"os"
)

func main() {

	log.SetLevel(log.TraceLevel)

	conf := common.NewConfig("")

	switch os.Args[1] {
	case "new":
		NewImage(conf, os.Args[2], os.Args[3])
	case "map":
		MapImage(conf, os.Args[2], os.Args[3])
	}

}
