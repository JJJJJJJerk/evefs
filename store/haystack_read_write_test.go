package store

import (
	"github.com/sirupsen/logrus"
	"hash/crc32"
	"io/ioutil"
	"testing"
)

func TestHaystack_AppendBlob_Write(t *testing.T) {
	data, err := ioutil.ReadFile("C:/Users/zhouqing1/Desktop/store.jpg")
	oCrc := crc32.ChecksumIEEE(data)

	if err != nil {
		t.Fatal(err)
	}

	hs := NewHaystack(9, 32*GB, "temp")
	offset, err := hs.WriteBlob(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(offset)
	logrus.Info(offset)

	img, err := hs.ReadBlob(int(offset), len(data))
	if err != nil {
		t.Fatal(err)
	}
	nCrc := crc32.ChecksumIEEE(img)

	if nCrc != oCrc {
		t.Errorf("file ")
	}

	ioutil.WriteFile("test4.jpg", img, 777)
}
