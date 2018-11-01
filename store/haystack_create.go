// +build !linux

package store

import (
	"os"
)

func (v *Haystack) createHayStackFile() error {
	if v.MaxSize > 0 {
		//logrus.Warn("windows do not support preallcate")
	}
	fileName := v.FileName()
	dataFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	v.dataFile = dataFile
	return nil
}
