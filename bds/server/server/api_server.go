package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetBackend(c *gin.Context) {
	c.JSON(http.StatusOK, s.Backend)
}
