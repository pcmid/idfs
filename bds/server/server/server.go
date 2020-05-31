package server

import (
	"github.com/gin-gonic/gin"
	"idfs/bds/common"
	"idfs/bds/common/image"
	"idfs/oss"
)

type Server struct {
	Addr    string
	Backend common.OssBackend

	imageMap map[string]*image.Image

	r *gin.Engine
}

func (s *Server) Init(conf *common.Config) {
	s.Addr = conf.Addr
	s.Backend = &oss.Server{Url: "192.168.0.3:9000/oss/"}

	s.imageMap = make(map[string]*image.Image)

	s.initRouter()

}

func (s *Server) initRouter() {
	s.r = gin.Default()

	serverG := s.r.Group("/server")
	{
		serverG.GET("/backend", s.GetBackend)
	}

	imageG := s.r.Group("/image")
	{
		imageG.PUT("/:name", s.CreateImage)
		imageG.GET("/:name", s.GetImage)
		imageG.PATCH("/:name", s.UpdateImage)
		imageG.DELETE("/:name", s.DeleteImage)

		blockG := imageG.Group("/:name/block")
		{
			blockG.GET("/:block", s.GetBlock)
		}
	}
}

func (s *Server) Run() {
	_ = s.r.Run(s.Addr)
}
