package main

import (
	log "github.com/sirupsen/logrus"
	"idfs/oss/node/config"
	"idfs/oss/node/server"
)

func main() {
	log.SetLevel(log.TraceLevel)

	conf := config.NewConfig("")
	s := &server.Server{}
	s.Init(conf)
	s.Run()
}
