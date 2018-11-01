package store

import (
	"github.com/dejavuzhou/evefs/pb"
	"github.com/gogo/protobuf/proto"
	"math/rand"
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

func (s *Store) WriteRpc(in *pb.NeedlePb) error {
	//choose a random haystack
	haystackIdx := rand.Intn(s.StackCount)
	hs := s.Stacks[haystackIdx]
	in.HaystackId = uint32(haystackIdx)
	//write file to haystack file
	err := hs.WriteNeedPb(in)
	//remove file bytes
	in.FileBytes = nil
	//save file info to store.levelDB
	valueBytes, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	err = s.DB.Put(IdToByets(in), valueBytes, nil)
	if err != nil {
		return err
	}
	return err
}
