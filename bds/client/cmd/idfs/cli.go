package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"idfs/bds/client/buse"
	"idfs/bds/client/device"
	"idfs/bds/common"
	"idfs/bds/common/image"
	"idfs/oss"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func NewImage(conf *common.Config, name string, size string) {
	requestUrl := fmt.Sprintf("http://%s/image/%s?size=%s", conf.Addr, name, size)
	resp, _ := resty.New().R().Put(requestUrl)
	img := image.EmptyImage()

	_ = json.Unmarshal(resp.Body(), img)
}

func MapImage(conf *common.Config, name string, nbd string) {

	requestUrl := fmt.Sprintf("http://%s/image/%s", conf.Addr, name)
	resp, _ := resty.New().R().Get(requestUrl)

	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("Failed to Map image '%s': %s", name, string(resp.Body()))
	}

	img := image.EmptyImage()

	_ = json.Unmarshal(resp.Body(), img)

	dev := &device.Device{
		Remote:  fmt.Sprintf("http://%s/image/%s", conf.Addr, img.Name),
		Name:    name,
		Backend: oss.NewServer(conf.Backend.Url),
		Img:     img,
	}

	go func() {
		tick := time.Tick(10 * time.Second)
		for {
			select {
			case <-tick:
				_ = dev.Flush()
			}
		}
	}()

	buseDev, err := buse.CreateDevice(nbd, dev.Img.Size, dev)
	if err != nil {
		log.Fatalf("Cannot create device: %s", err)
		//os.Exit(1)
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		if err := buseDev.Connect(); err != nil {
			log.Printf("Buse device stopped with error: %s", err)
		} else {
			log.Println("Buse device stopped gracefully.")
		}
	}()
	<-sig
	// Received SIGTERM, cleanup
	fmt.Println("SIGINT, disconnecting...")
	buseDev.Disconnect()
}
