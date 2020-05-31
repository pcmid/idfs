package storage

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Disk struct {
	//Name       string
	MountPoint string
}

func (d *Disk) Write(filename string, content []byte) error {
	if err := ioutil.WriteFile(d.MountPoint+"/"+filename, content, 0600); err != nil {
		log.Errorf("Failed to write object: %server", err)
		return err
	}
	return nil
}

func (d *Disk) Read(filename string) []byte {
	if content, err := ioutil.ReadFile(d.MountPoint + "/" + filename); err != nil {
		log.Errorf("Failed to read object: %server", err)
		return nil
	} else {
		return content
	}
}

func (d *Disk) Delete(filename string) {
	if err := os.Remove(d.MountPoint + "/" + filename); err != nil {
		log.Errorf("Failed to remove object: %server", err)
	}
}
