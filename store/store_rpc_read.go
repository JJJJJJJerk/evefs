package store

import (
	"github.com/dejavuzhou/evefs/pb"
	"github.com/gogo/protobuf/proto"
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

func (s *Store) ReadRpc(in *pb.NeedlePb) error {
	idKey := IdToByets(in)
	buf, err := s.DB.Get(idKey, nil)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(buf, in)
	if err != nil {
		return err
	}
	hs := s.Stacks[in.HaystackId]
	
	err = hs.ReadFileBytes(in)
	return err
}
