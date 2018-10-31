package store

import (
	"errors"
	"fmt"
	"github.com/dejavuzhou/evefs/pb"
	"github.com/dejavuzhou/evefs/snowflake"
	"mime/multipart"
)

/**

		-----------------------------
		|-headerCrc32 8 bytes  first 4 byes is crc32Uint 5th byte is pad-length, 6th 7th 8th bytes are reserved
		|-magicHeader 8  bytes
 needle-|
		|-FileBytes + padding bytes 8 * n bytes
		|-magicFooter 8  bytes
		-----------------------------

*/

// writeNeedleBytes append a blob to end of the FileBytes file, used in replication
func (hs *Haystack) writeNeedleBytes(n *pb.NeedlePb) (err error) {
	offset, err := hs.dataFile.Seek(0, 2)
	if err != nil {
		return errors.New("failed to seek the end of file")
	}
	dataPaddingSize := NeedlePaddingSize - int64(n.FileSize)%NeedlePaddingSize
	needleLength := int64(n.FileSize) + NeedlePaddingSize*3 + dataPaddingSize
	//TODO:: 不知道这个bytes size 是否只是正确的
	if offset+needleLength > int64(hs.MaxSize) {
		err := errors.New("volume Storage is full. can not append file into it")
		return err
	}
	//ensure file writing starting from aligned positions
	if offset%NeedlePaddingSize != 0 {
		offset = offset + (NeedlePaddingSize - offset%NeedlePaddingSize)
		if offset, err = hs.dataFile.Seek(offset, 0); err != nil {
			msg := fmt.Sprintf("failed to align in datafile %s", hs.dataFile.Name())
			return errors.New(msg)
		}
	}
	//create needle bytes
	needleByets, err := createNeedleBytes(n)
	if err != nil {
		return err
	}
	hs.fileLock.Lock()
	_, err = hs.dataFile.Write(needleByets)
	hs.fileLock.Unlock()
	if err == nil {
		n.Offset = uint32(offset / NeedlePaddingSize)
	}
	return
}

func (hs *Haystack) WriteFileHeader(fh *multipart.FileHeader) (*pb.NeedlePb, error) {
	node, err := snowflake.NewNode(int64(hs.Id))
	if err != nil {
		return nil, err
	}
	needleId := node.Generate()
	
	n, err := NewNeedle(int64(needleId), uint8(hs.Id), fh)
	if err != nil {
		return nil, err
	}
	if err := hs.writeNeedleBytes(n); err != nil {
		return nil, err
	}
	return n, err
}
func (hs *Haystack) WriteNeedPb(n *pb.NeedlePb) (error) {
	
	if err := hs.writeNeedleBytes(n); err != nil {
		return err
	}
	return nil
}