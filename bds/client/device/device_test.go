package device

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"idfs/bds/client/buse"
	"idfs/oss"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestDevice_ReadWrite(t *testing.T) {

	_ = uint(1024 * 1024 * 512) // 512M
	device := &Device{
		Remote:  "http://127.0.0.1:8000/image/test",
		Name:    "test",
		Backend: &oss.Server{Url: "http://127.0.0.1:9000/oss/"},
		Img:     nil,
	}

	device.Fetch()

	device.WriteAt([]byte("abcdef\n"), 0)

	data := make([]byte, 7)
	device.ReadAt(data, 0)

	device.Flush()

	fmt.Printf("%s", string(data))
}

func TestDevice_Flush(t *testing.T) {
	//nbdDev := os.Args[1]

	log.SetLevel(log.TraceLevel)

	//size := uint(1024 * 1024 * 512) // 512M
	deviceExp := &Device{
		Remote:  "http://192.168.0.3:8000/image/test",
		Name:    "test",
		Backend: &oss.Server{Url: "http://192.168.0.3:9000/oss/"},
		Img:     nil,
	}

	go func() {
		tick := time.Tick(1 * time.Second)
		for {
			select {
			case <-tick:
				_ = deviceExp.Flush()
			}
		}
	}()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
	}()

	deviceExp.Fetch()

	device, err := buse.CreateDevice("/dev/nbd0", deviceExp.Img.Size, deviceExp)
	if err != nil {
		log.Fatalf("Cannot create device: %s", err)
		//os.Exit(1)
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		if err := device.Connect(); err != nil {
			log.Printf("Buse device stopped with error: %s", err)
		} else {
			log.Println("Buse device stopped gracefully.")
		}
	}()

	time.Sleep(20 * time.Second)

	//<-sig
	// Received SIGTERM, cleanup
	fmt.Println("disconnecting...")
	device.Disconnect()
}
