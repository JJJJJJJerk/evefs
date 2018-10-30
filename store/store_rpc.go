package store

import "github.com/dejavuzhou/evefs/pb"

func (s *Store) RpcWrite(in *pb.PutData) (*pb.NeedlePb, error) {
	return nil, nil
}
func (s *Store) RpcRead(in *pb.NeedlePb) (*pb.PutData, error) {
	return nil, nil
	
}

func (s *Store) RpcStatus(in *pb.Node) (*pb.NodeStatus, error) {
	return nil, nil
}
