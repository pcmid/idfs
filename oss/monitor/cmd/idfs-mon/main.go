package main

import (
	"idfs/oss/monitor/config"
	"idfs/oss/monitor/server"
)

func main() {
	c := config.NewConfig("")
	s := &server.Server{}
	s.Init(c)
	s.Run()
}
