package server

import (
	"github.com/gin-gonic/gin"
	"idfs/oss/monitor/config"
	"idfs/oss/monitor/ring"
)

type Server struct {
	Addr string
	Ring *ring.Ring

	r *gin.Engine
}

func (s *Server) Init(conf *config.Config) {
	s.Addr = conf.Addr

	s.Ring = ring.NewRing(conf.RingPath)

	s.initRouter()
}

func (s *Server) initRouter() {
	s.r = gin.Default()
	backendG := s.r.Group("/oss")
	{
		backendG.GET("/:item", s.GetBackend)
	}
}

func (s *Server) Run() {
	_ = s.r.Run(s.Addr)
}
