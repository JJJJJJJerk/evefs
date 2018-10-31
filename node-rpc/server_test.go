package main

import (
	"context"
	"github.com/dejavuzhou/evefs/pb"
	"google.golang.org/grpc"
	"hash/crc32"
	"io/ioutil"
	"log"
	"testing"
)

const (
	address = "localhost:12345"
)

func TestRpcServer_WriteFile(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRpcServiceClient(conn)
	
	// Contact the server and print out its response.
	
	n := createWriteNeedlePb()
	r, err := c.WriteFile(context.Background(), n)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r)
	
}

func createWriteNeedlePb() (*pb.NeedlePb) {
	fileData, _ := ioutil.ReadFile("C:/Users/zhouqing1/Desktop")
	n := &pb.NeedlePb{}
	n.Id = 12345678
	n.HaystackId = 1
	n.MimeType = "image/jpg"
	n.CheckSum = crc32.ChecksumIEEE(fileData)
	n.FileSize = uint32(len(fileData))
	n.FileBytes = fileData
	n.Offset = 0
	n.Flags = 0
	return n
}
