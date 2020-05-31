package oss

import (
	"crypto/md5"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type MonitorResp struct {
	Url string `json:"url"`
}

type Server struct {
	Url string
}

func NewServer(url string) *Server {
	return &Server{Url: url}
}

func (s *Server) Get(item string) ([]byte, error) {

	resp, err := resty.New().R().Get(s.Url + item)

	if err != nil {
		log.Errorf("Get item: failed to get item url: %s", err)
		return nil, err
	}

	mResp := &MonitorResp{}
	err = json.Unmarshal(resp.Body(), mResp)

	dataResp, e := resty.New().R().Get(mResp.Url)

	if e != nil {
		log.Errorf("Failed to get item '%s': %s", item, e)
		return nil, e
	}

	data := dataResp.Body()

	log.Debugf("Get item '%s' from '%s', size: %dM", item, s.Url+item, len(data)/1024/1024)

	md5Ctx := md5.New()
	md5Ctx.Write(data)
	log.Tracef("md5sum: %x", md5.Sum(data))
	return data, nil
}

func (s *Server) Put(item string, data []byte) {

	resp, err := resty.New().R().Get(s.Url + item)

	if err != nil {
		log.Errorf("Put item: failed to get item url: %s", err)
		return
	}

	mResp := &MonitorResp{}
	_ = json.Unmarshal(resp.Body(), mResp)

	itemUrl := mResp.Url

	log.Debugf("Put item '%s' to '%s', size: %dM", item, itemUrl, len(data)/1024/1024)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	log.Tracef("md5sum: %x", md5.Sum(data))

	_, err = resty.New().R().
		SetBody(data).
		SetContentLength(true).
		Put(itemUrl)

	if err != nil {
		log.Errorf("%s", err)
	}
}

func (s *Server) Delete(item string) {

}
