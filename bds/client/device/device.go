package device

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"idfs/bds/common"
	"idfs/bds/common/image"
	"runtime"
	"runtime/debug"
	"time"
)

type Device struct {
	Remote string
	Name   string

	Backend common.OssBackend
	Img     *image.Image
}

func (d *Device) Fetch() error {
	log.Debug("Fetch metadata")

	var imageData []byte

	if resp, err := resty.New().R().Get(d.Remote); err != nil {
		return err
	} else {
		imageData = resp.Body()
	}

	if d.Img == nil {
		d.Img = image.EmptyImage()
		_ = json.Unmarshal(imageData, d.Img)
		return nil
	}

	// TODO
	imageRemote := image.EmptyImage()
	_ = json.Unmarshal(imageData, imageRemote)

	return nil

}

func (d *Device) pullBlock(block *image.Block) {
	log.Tracef("Pull block %s", block.ID.String())

	if !block.Created {
		log.Tracef("Create New block: %s", block.ID.String())
		block.Cached(nil)
		block.Created = true
		block.LastUpdate = time.Now()
		return
	}

	if block.Cache == nil {
		log.Tracef("Cache block: %s", block.ID.String())
		cache, err := d.Backend.Get(block.ID.String())
		if err != nil {
			log.Error(err)
			block.Created = false
			d.pullBlock(block)
		}
		block.Cached(cache)
	}
}

func (d *Device) ReadAt(p []byte, off uint64) error {

	if d.Img == nil {
		if err := d.Fetch(); err != nil {
			return err
		}
	}

	block, pos := d.Img.BlockAt(off)

	d.pullBlock(block)

	n := block.ReadAt(p, pos)

	if n == 0 {
		return errors.New("nil data for block")
	}

	if n < uint64(len(p)) {
		return d.ReadAt(p[n:], off+n)
	}
	return nil
}

func (d *Device) WriteAt(p []byte, off uint64) error {
	if d.Img == nil {
		if err := d.Fetch(); err != nil {
			return err
		}
	}

	block, pos := d.Img.BlockAt(off)
	//log.Tracef("Write Block %s", block.ID.String())

	d.pullBlock(block)

	n := block.WriteAt(p, pos)

	if n == 0 {
		return errors.New("nil data for block")
	}

	if n < uint64(len(p)) {
		return d.WriteAt(p[n:], off+n)
	}

	return nil
}

func (d *Device) Disconnect() {
	_ = d.Flush()
	log.Info("Disconnected")
}

func (d *Device) Flush() error {
	log.Debug("Flush device")

	if d.Img == nil {
		return nil
	}

	for _, b := range d.Img.Blocks {
		b.Lock()
		if b.Cache == nil || len(b.Cache) == 0 {
			b.Unlock()
			//log.Tracef("Null data: %s", b.ID.String())
			continue
		}

		d.Backend.Put(b.ID.String(), b.Cache)
		b.Cache = nil
		b.Unlock()
		//log.Tracef("%s size: %dM", b.ID.String(), len(b.Cache)/1024/1024)
	}

	_, err := resty.New().R().
		SetBody(d.Img).
		SetJSONEscapeHTML(true).
		SetContentLength(true).
		Patch(d.Remote)

	runtime.GC()
	debug.FreeOSMemory()

	return err
}

func (d *Device) Trim(off uint64, length uint) error {
	log.Debugf("Trim at: %d, length: %d", off, length)
	//trimData := make([]byte, length)
	//d.WriteAt(trimData, off)
	return nil
}
