package store

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"hash/crc32"
	"io/ioutil"
	"mime/multipart"
	"strconv"
)

type Needle struct {
	Id      int64  `comment:"random number to mitigate brute force lookups"`
	StackId uint8  `comment:"random number to mitigate brute force lookups"`
	Offset  uint32 `comment:"sum of DataSize,Data,NameSize,GetDataDirPath,MimeSize,Mime"`
	Flags   uint8  `comment:"boolean flags"`          //version2
	Name    string `comment:"maximum 256 characters"` //version2

	Size uint32 `comment:"sum of DataSize,Data,NameSize,GetDataDirPath,MimeSize,Mime"`
	Mime string `comment:"maximum 256 characters"` //version2
	//Checksum   uint32 `comment:"CRC32 to check integrity"`
	data     []byte `comment:"The actual file data"`
	checkSum uint32
}

func (n *Needle) IdToByets() []byte {
	numbers := strconv.FormatInt(n.Id, 10)
	return []byte(numbers)
}

func NewNeedle(needleId int64, stackId uint8, file *multipart.FileHeader) (needle *Needle, err error) {
	needle = &Needle{}
	needle.Id = needleId
	needle.Flags = 0
	needle.StackId = stackId
	needle.Name = file.Filename
	needle.Size = uint32(file.Size)
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
	needle.data = data
	needle.checkSum = crc32.ChecksumIEEE(data)
	return needle, nil
}

func (n *Needle) LevelDbCrc32Key() (crcBytes []byte) {
	if n.Size < 1 || n.checkSum < 0 {
		err := errors.New("needle has not initialized fully, data, checkSum, and size are required")
		logrus.Fatal(err)
	}
	crc32Size := fmt.Sprint("%x,%x", n.checkSum, n.Size)
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
func (n *Needle) CreateWriteBytes() (blob []byte, err error) {
	if n.Size == 0 || len(n.data) == int(n.Size) || n.checkSum == 0 {
		return nil, errors.New("needle has not initialized fully, data,checkSum, and size are required")
	}
	// the header is 8 byte
	header := make([]byte, NeedlePaddingSize, NeedlePaddingSize)
	//write first 4 bytes with crc32 return value
	//the reset 4 bytes are unused
	uint32toBytes(header, n.checkSum)
	//create paddingBytes bytes to align
	padNum := int(NeedlePaddingSize) - len(n.data)%int(NeedlePaddingSize)
	paddingBytes := make([]byte, padNum, padNum)
	//write bytes in following order header(crc32) + data + paddingBytes
	buf := bytes.Buffer{}
	_, err = buf.Write(header)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(n.data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(paddingBytes)
	if err != nil {
		return nil, err
	}
	n.data = nil
	blobBinary := buf.Bytes()
	return blobBinary, err
}

func (n *Needle) paresData(b []byte) error {
	//crs32 checksum returning uint32 only uses 4 byte, the rest 4 bytes are unused
	headerCrc32bytes := b[0:4]
	checkSum := bytesToUint32(headerCrc32bytes)
	//get file data bytes
	toIdx := n.Size + uint32(NeedlePaddingSize)
	fromIdx := NeedlePaddingSize * 1
	data := b[fromIdx:toIdx]
	//check file data's CRC32
	if crc32.ChecksumIEEE(data) != checkSum {
		err := errors.New("needle parse error. volume file is broken")
		return err
	}
	n.data = data
	return nil
}

//func NewNeedleFileHeader(file *multipart.FileHeader) (n *Needle, err error) {
//
//
//	n = Needle{}
//	return nil, err
//}
