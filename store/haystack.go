package store

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"sync"
)

type Haystack struct {
	Id       uint32
	dataFile *os.File
	fileLock sync.RWMutex
	MaxSize  ByteSize
	barnDir  string
}

func NewHaystack(id int, maxSize ByteSize, barnDir string) *Haystack {
	h := Haystack{}
	h.barnDir = barnDir
	h.fileLock = sync.RWMutex{}
	h.Id = uint32(id)
	if uint32(maxSize) == 0 {
		h.MaxSize = math.MaxUint32
	} else {
		h.MaxSize = ByteSize(maxSize)
	}
	if h.createHayStackFile() != nil {
		logrus.Fatalln("could not create Haystack dataFile")
	}
	return &h
}

func (h *Haystack) FileName() string {
	return fmt.Sprintf("%s/%d.data", h.barnDir, h.Id)
}
