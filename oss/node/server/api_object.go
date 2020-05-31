package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) PutObject(c *gin.Context) {
	disk := c.Param("disk")
	objectName := c.Param("obj")

	dataLen, _ := strconv.ParseUint(c.Request.Header.Get("Content-Length"), 10, 64)

	data := make([]byte, dataLen)
	readLen := uint64(0)

	for readLen < dataLen {
		n, err := c.Request.Body.Read(data[readLen:])
		if err != nil && err != io.EOF {
			log.Errorf("Failed to put object '%s'", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to put object '%s'", objectName),
			})
			return
		}
		readLen += uint64(n)
	}

	log.Tracef("Read object size: %d", readLen)

	err := s.GetDisk(disk).Write(objectName, data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to put object '%server'", objectName),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Put object '%server'", objectName),
	})
}

func (s *Server) GetObject(c *gin.Context) {
	disk := c.Param("disk")
	objectName := c.Param("obj")

	data := s.GetDisk(disk).Read(objectName)

	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (s *Server) DeleteObject(c *gin.Context) {
	disk := c.Param("disk")
	objectName := c.Param("obj")

	s.GetDisk(disk).Delete(objectName)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Delete object '%server'", objectName),
	})

}
