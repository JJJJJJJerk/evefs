package store

import (
	"errors"
	"github.com/dejavuzhou/spookyfs/snowflake"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
)

// Destroy removes everything related to this volume
func (hs *Haystack) Destroy() (err error) {
	err = os.Remove(hs.dataFile.Name())
	if err != nil {
		return
	}
	os.Remove(hs.FileName() + ".cpd")
	os.Remove(hs.FileName() + ".cpx")
	os.Remove(hs.FileName() + ".ldb")
	os.Remove(hs.FileName() + ".bdb")
	return
}

// WriteBlob append a blob to end of the data file, used in replication
func (hs *Haystack) WriteBlob(b []byte) (offset int64, err error) {
	if offset, err = hs.dataFile.Seek(0, 2); err != nil {
		logrus.WithError(err).Error("failed to seek the end of file")
		return
	}
	if int(offset)+len(b) > int(hs.MaxSize) {
		err := errors.New("volume Storage is full. can not append file into it")
		logrus.WithError(err).Errorf("volumen %d is full", hs.Id)
		return 0, err
	}
	//ensure file writing starting from aligned positions
	if offset%NeedlePaddingSize != 0 {
		offset = offset + (NeedlePaddingSize - offset%NeedlePaddingSize)
		if offset, err = hs.dataFile.Seek(offset, 0); err != nil {
			logrus.WithError(err).Errorf("failed to align in datafile %s", hs.dataFile.Name())
			return
		}
	}
	if offset%NeedlePaddingSize != 0 {
		logrus.Fatalln("calculate file offset wrong")
	} else {
		offset = offset / NeedlePaddingSize
	}
	hs.fileLock.Lock()
	defer hs.fileLock.Unlock()
	_, err = hs.dataFile.Write(b)
	return
}
func (hs *Haystack) ReadBlob(offset, size int) (dataSlice []byte, err error) {
	osOffset := int64(offset) * NeedlePaddingSize
	if int(osOffset)+size >= int(hs.MaxSize) {
		err = errors.New("offest is out range of volume")
		logrus.WithError(err).Error("offest is out range of volume. id is ", hs.Id)
	}
	dataSlice = make([]byte, int(size))

	hs.fileLock.RLock()
	hs.fileLock.RUnlock()
	_, err = hs.dataFile.ReadAt(dataSlice, osOffset)
	return dataSlice, err
}

func (hs *Haystack) PutFileHead(fh *multipart.FileHeader) (*Needle, error) {
	//generate snowflake id
	node, err := snowflake.NewNode(int64(hs.Id))
	if err != nil {
		return nil, err
	}
	needleId := node.Generate()

	//create needle with basic info
	needle, err := NewNeedle(int64(needleId), uint8(hs.Id), fh)
	if err != nil {
		return nil, err
	}
	//create write bytes blob
	writeBytes, err := needle.CreateWriteBytes()
	if err != nil {
		return nil, err
	}
	//write byte blob to haystack file
	//return offset and gives it to need
	offset, err := hs.WriteBlob(writeBytes)
	needle.Offset = uint32(offset)
	return needle, err
}
