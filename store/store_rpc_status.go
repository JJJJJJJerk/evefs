package store

import (
	"github.com/dejavuzhou/evefs/pb"
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

func (s *Store) StatusRpc(in *pb.NeedlePb) error {
	return nil
}
