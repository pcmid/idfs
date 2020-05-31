package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"hash/adler32"
	"idfs/oss/node/config"
	"idfs/oss/node/storage"
	"strconv"
)

type Server struct {
	Addr string

	Disks map[uint32]*storage.Disk

	r *gin.Engine
}

func (s *Server) Init(conf *config.Config) {
	s.Addr = conf.Addr

	s.Disks = make(map[uint32]*storage.Disk, len(conf.Disks))

	for _, disk := range conf.Disks {
		id := adler32.Checksum([]byte(disk))
		log.Debugf("disk '%s' as '%d'", disk, id)
		s.Disks[id] = &storage.Disk{MountPoint: disk}
	}

	s.initRouter()
}

func (s *Server) initRouter() {
	s.r = gin.Default()

	objG := s.r.Group("/obj")
	{
		objG.GET("/:disk/:obj", s.GetObject)
		objG.PUT("/:disk/:obj", s.PutObject)
		objG.DELETE("/:disk/:obj", s.DeleteObject)
	}
}

func (s *Server) Run() {
	_ = s.r.Run(s.Addr)
}

func (s *Server) GetDisk(disk string) *storage.Disk {
	id, _ := strconv.ParseUint(disk, 10, 32)
	return s.Disks[uint32(id)]
}
