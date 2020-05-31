package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"idfs/bds/common/image"
	"net/http"
	"strconv"
)

type Response struct {
	Error_ *Error       `json:"error"`
	Image_ *image.Image `json:"image"`
}

type Error struct {
	Massage string `json:"massage"`
}

func (s *Server) GetImage(c *gin.Context) {

	imageName := c.Param("name")

	if img, ok := s.imageMap[imageName]; ok {
		c.JSON(http.StatusOK, img)
	} else {
		c.JSON(http.StatusNotFound, Error{
			Massage: fmt.Sprintf("image '%s' not found", imageName),
		})
	}

	return
}

func (s *Server) CreateImage(c *gin.Context) {
	imageName := c.Param("name")
	imageSize := uint64(0)

	if _, ok := s.imageMap[imageName]; ok {
		c.JSON(http.StatusConflict, Error{Massage: fmt.Sprintf("image '%server' has already exist", imageName)})
		return
	}

	if size, ok := c.GetQuery("size"); ok == false {
		c.JSON(http.StatusBadRequest, Error{
			Massage: "argument size not found",
		})
		return
	} else {
		imageSize, _ = strconv.ParseUint(size, 10, 64)
	}

	img := image.NewImage(imageName, imageSize)

	s.imageMap[imageName] = img

	c.JSON(http.StatusOK, Response{
		Error_: nil,
	})
}

func (s *Server) UpdateImage(c *gin.Context) {
	imageName := c.Param("name")

	image := s.imageMap[imageName]

	dataLen := c.Request.ContentLength
	imageData := make([]byte, dataLen)

	readLen := 0

	for int64(readLen) < dataLen {
		n, _ := c.Request.Body.Read(imageData[readLen:])
		readLen += n
	}

	_ = json.Unmarshal(imageData, image)

	c.JSON(http.StatusOK, gin.H{
		"message": "Update image:" + imageName,
	})
}

func (s *Server) DeleteImage(c *gin.Context) {
	imageName := c.Param("name")

	if _, ok := s.imageMap[imageName]; ok {
		delete(s.imageMap, imageName)

		c.JSON(http.StatusOK, Response{
			Error_: nil,
		})
	} else {
		c.JSON(http.StatusNotFound, Error{
			Massage: fmt.Sprintf("image '%server' not found", imageName),
		})
	}
}
