package store

import (
	"io/ioutil"
	"testing"
)

func TestNewBarn(t *testing.T) {
	NewStore("127.0.0.1:8787", "temp", 8)
}

func TestBarn_PutFile(t *testing.T) {
	barn := NewStore("127.0.0.1:8787", "temp", 8)
	data, err := ioutil.ReadFile("C:/Users/zhouqing1/Desktop/store.jpg")
	if err != nil {
		t.Error(err)
	}
	barn.PutFile(data)
}
