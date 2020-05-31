package main

import (
	"idfs/bds/common"
	"idfs/bds/server/server"
)

func main() {
	conf := common.NewConfig("")
	s := &server.Server{}
	s.Init(conf)
	s.Run()
}