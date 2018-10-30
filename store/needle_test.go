package store

import (
	"io/ioutil"
	"testing"
)

func TestBarn_NeedleParseData(t *testing.T) {
	needle := Needle{}
	needle.Id = 777
	needle.StackId = 1
	needle.Flags = 2
	needle.Offset = 1
	needle.Name = "aweomse.jpg"
	data, _ := ioutil.ReadFile("C:/Users/zhouqing1/Desktop/store.jpg")
	needle.FileBytes = data
	needle.Size = uint32(len(data))

	data, err := needle.CreateWriteBytes()
	if err != nil {
		t.Error(err)
	}
	needle.paresData(data)

}
