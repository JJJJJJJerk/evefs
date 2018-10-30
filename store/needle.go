package store

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dejavuzhou/evefs/pb"
	"github.com/sirupsen/logrus"
	"hash/crc32"
	"io/ioutil"
	"mime/multipart"
	"strconv"
)

var (
	// use magic header/footer 8 bytes for repairing FileBytes
	magicHeader = []byte("felixeve")
	magicFooter = []byte("felixout")
)

//needle

//type Needle struct {
//	Id      int64
//	StackId uint8
//	Offset  uint32
//	Flags   uint8
//	Name    string
//
//	Size uint32
//	Mime string
//	//Checksum   uint32 `comment:"CRC32 to check integrity"`
//	FileBytes []byte `json:"-"`
//	CheckSum  uint32
//}

type Needle struct {
	pb.NeedlePb
	FileBytes []byte `json:"-"`
}

func (n *Needle) IdToByets() []byte {
	numbers := strconv.FormatUint(n.Id, 10)
	return []byte(numbers)
}

func NewNeedle(needleId int64, stackId uint8, file *multipart.FileHeader) (needle *Needle, err error) {
	needle = &Needle{}
	needle.Id = uint64(needleId)
	needle.Flags = 0
	needle.StackId = uint32(stackId)
	needle.Name = file.Filename
	needle.Size = uint32(file.Size)
	//TODO 这个地方有可能会panic
	heads := file.Header["Content-Type"]
	needle.Mime = heads[0]

	tempFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(tempFile)
	if err != nil {
		return nil, err
	}
	needle.FileBytes = data
	needle.CheckSum = crc32.ChecksumIEEE(data)
	return needle, nil
}

func (n *Needle) LevelDbCrc32Key() (crcBytes []byte) {
	if n.Size < 1 || n.CheckSum < 0 {
		err := errors.New("needle has not initialized fully, FileBytes, CheckSum, and size are required")
		logrus.Fatal(err)
	}
	crc32Size := fmt.Sprint("%x,%x", n.CheckSum, n.Size)
	return []byte(crc32Size)
}

func uint32toBytes(b []byte, v uint32) {
	for i := uint(0); i < 4; i++ {
		b[3-i] = byte(v >> (i * 8))
	}
}
func bytesToUint32(b []byte) (v uint32) {
	length := uint(len(b))
	for i := uint(0); i < length-1; i++ {
		v += uint32(b[i])
		v <<= 8
	}
	v += uint32(b[length-1])
	return
}

/**
		-----------------------------
		|-headerCrc32 8 bytes  first 4 byes is crc32Uint 5th byte is pad-length, 6th 7th 8th bytes are reserved
		|-magicHeader 8  bytes
 needle-|
		|-FileBytes + padding bytes 8 * n bytes
		|-magicFooter 8  bytes
		-----------------------------
*/

//func NewNeedleFileHeader(file *multipart.FileHeader) (n *Needle, err error) {
//
//
//	n = Needle{}
//	return nil, err
//}

func (n *Needle) createNeedleBytes() (needleBytes []byte, err error) {
	if n.Size == 0 || len(n.FileBytes) != int(n.Size) || n.CheckSum == 0 {
		return nil, errors.New("needle has not initialized fully, FileBytes,CheckSum, and size are required")
	}
	// the headerCrc32 is 8 byte
	headerCrc32 := make([]byte, NeedlePaddingSize, NeedlePaddingSize)
	//write first 4 bytes with crc32 return value
	//the reset 4 bytes are unused
	uint32toBytes(headerCrc32, n.CheckSum)
	//create paddingBytes bytes to align
	padNum := int(NeedlePaddingSize) - len(n.FileBytes)%int(NeedlePaddingSize)
	
	//set the headerCrc32 5th byte padNum
	headerCrc32[4] = byte(padNum)
	//paddingBytes
	paddingBytes := make([]byte, padNum, padNum)
	//write bytes in following order headerCrc32(crc32) + FileBytes + paddingBytes
	buf := bytes.Buffer{}
	_, err = buf.Write(headerCrc32)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(magicHeader)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(n.FileBytes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(paddingBytes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(magicFooter)
	if err != nil {
		return nil, err
	}
	n.FileBytes = nil
	needleBytes = buf.Bytes()
	return needleBytes, nil
}
