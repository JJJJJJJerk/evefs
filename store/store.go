package store

import (
	"fmt"
	"github.com/dejavuzhou/evefs/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"math/rand"
	"mime/multipart"
	"strings"
	"sync"
)

type Store struct {
	IpPort       string
	Stacks       []*Haystack
	StackCount   int
	StackMaxSize ByteSize
	DB           *leveldb.DB
	dataDir      string
	idLock       sync.Mutex
}

func (s *Store) Close() {
	s.DB.Close()
}

func (s *Store) GetDataDirPath() string {
	if s.IpPort == "" {
		logrus.Fatalln("Store's IpPort property has not set")
	}
	if s.StackCount < 1 {
		logrus.Fatalln("Store's StackCount property has not set or fatal error")
	}
	if s.dataDir == "" {
		s.dataDir = "."
	}

	temp := strings.Replace(s.IpPort, ".", "_", -1)
	temp = strings.Replace(temp, ":", "p", -1)
	return fmt.Sprintf("%s/%sc%d", s.dataDir, temp, s.StackCount)
}
func (s *Store) createLevelDb() {
	levelDbName := fmt.Sprintf("%s/leveldb", s.GetDataDirPath())
	db, err := leveldb.OpenFile(levelDbName, nil)
	if err != nil {
		logrus.WithError(err).Fatal("new Store could not create level DB")
	}
	s.DB = db
}



func NewStore(ipPort string, dataDir string, stackCount int) *Store {
	b := Store{}
	b.IpPort = ipPort
	b.StackCount = stackCount
	b.StackMaxSize = StackMaxSize
	b.dataDir = dataDir

	b.createLevelDb()
	for i := 0; i < stackCount; i++ {
		hs := NewHaystack(i, StackMaxSize, b.GetDataDirPath())
		b.Stacks = append(b.Stacks, hs)
	}
	return &b
}

func (s *Store) PutFile(file *multipart.FileHeader) (*pb.NeedlePb, error) {
	//get a random stackId
	stackId := rand.Intn(s.StackCount)
	//write binary
	hs := s.Stacks[stackId]
	needle, err := hs.WriteFileHeader(file)
	if err != nil {
		return nil, err
	}
	//write needle FileBytes to level db
	jsonData, err := proto.Marshal(needle)
	//TODO:: 可以检查文件是否存在直接返回相同的offset
	err = s.DB.Put(IdToByets(needle), jsonData, nil)
	if err != nil {
		return nil, err
	}
	return needle, nil
}

func (s *Store) DeleteFile(id int) {

}
func (s *Store) GetOneFile(id int) []byte {
	return nil
}
func (s *Store) GetAllFile() {

}

func (s *Store) getNeedleFromDb(id []byte) (needle *pb.NeedlePb, err error) {
	jsonData, err := s.DB.Get(id, nil)
	if err != nil {
		return nil, err
	}
	var n = &pb.NeedlePb{}
	proto.Unmarshal(jsonData, n)
	return n, nil
}

func (s *Store) GetFile(id []byte) (n *pb.NeedlePb, err error) {
	n, err = s.getNeedleFromDb(id)
	if err != nil {
		return nil, err
	}
	hs := s.Stacks[n.HaystackId]
	err = hs.ReadFileBytes(n)
	if err != nil {
		return n, err
	}
	
	return n, nil
}
