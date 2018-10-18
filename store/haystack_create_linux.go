// +build linux

package store

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func (v *Haystack) createHayStackFile() error {
	fileName := v.FileName()
	dataFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if v.MaxSize > 0 {
		syscall.Fallocate(int(file.Fd()), 1, 0, int64(v.MaxSize))
		logrus.Infof("Preallocated %d bytes disk space for %s", preallocate, fileName)
	}
	v.dataFile = dataFile
	return nil
}
