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

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) WriteFile(ctx context.Context, in *pb.PutData) (*pb.NeedlePb, error) {
	return nil, nil
}
func (s *server) ReadFile(ctx context.Context, in *pb.NeedlePb) (*pb.PutData, error) {
	//
	return nil, nil
}
func (s *server) Status(ctx context.Context, in *pb.Node) (*pb.NodeStatus, error) {
	return nil, nil
}
