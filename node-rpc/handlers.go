package main

import (
	"github.com/dejavuzhou/evefs/pb"
	"github.com/dejavuzhou/evefs/store"
)
import "context"

var nodeStore *store.Store

func init() {
	//TODO::read address and port from config
	nodeStore = store.NewStore("127.0.0.1:1212", "temp", 8)
	
}

// rpcServer is used to implement helloworld.GreeterServer.
type rpcServer struct {
	AllowIp string
}

func (s *rpcServer) WriteFile(ctx context.Context, in *pb.NeedlePb) (*pb.NeedlePb, error) {
	hs := nodeStore.Stacks[in.HaystackId]
	hs.WriteNeedPb(in)
	return in, nil
}
func (s *rpcServer) ReadFile(ctx context.Context, in *pb.NeedlePb) (*pb.NeedlePb, error) {
	
	return in, nil
}
func (s *rpcServer) Status(ctx context.Context, in *pb.Node) (*pb.NodeStatus, error) {

	return nil, nil
}

