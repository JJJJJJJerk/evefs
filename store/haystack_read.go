package store

import (
	"errors"
	"github.com/dejavuzhou/evefs/pb"
	"hash/crc32"
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
func (hs *Haystack) ReadFileBytes(n *pb.NeedlePb) error {
	
	osOffset := int64(n.Offset) * NeedlePaddingSize
	padNum := NeedlePaddingSize - int64(n.FileSize)%NeedlePaddingSize
	needleBytes := make([]byte, int64(n.FileSize)+NeedlePaddingSize*3+padNum)
	
	hs.fileLock.RLock()
	_, err := hs.dataFile.ReadAt(needleBytes, osOffset)
	hs.fileLock.RUnlock()
	if err != nil {
		return err
	}
	if byte(padNum) != needleBytes[4] {
		return errors.New("check needle FileBytes pad number failed")
	}
	crc32Bytes := needleBytes[0:4]
	crcValue := bytesToUint32(crc32Bytes)
	
	if crcValue != n.CheckSum {
		return errors.New("haystack file crc32 bytes are wrong")
	}
	fromIdx := 2 * NeedlePaddingSize
	toIdx := fromIdx + int64(n.FileSize)
	fileData := needleBytes[fromIdx:toIdx]
	
	if crc32.ChecksumIEEE(fileData) != crcValue {
		return errors.New("file FileBytes crc32 check failed")
	}
	n.FileBytes = fileData
	n.CheckSum = crcValue
	return err
}
