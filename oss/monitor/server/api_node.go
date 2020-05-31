package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetBackend(c *gin.Context) {
	item := c.Param("item")

	node := s.Ring.Get(item)

	c.JSON(http.StatusOK, gin.H{
		"url": node.Url() + item,
	})

}
