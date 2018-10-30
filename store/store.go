package store

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"math/rand"
	"mime/multipart"
	"strings"
	"sync"
)

type Store struct {
	IpPort                   string
	Stacks                   []*Haystack
	StackCount               int
	StackMaxSize             ByteSize
	Db                       *leveldb.DB
	NeedleMaxAutoIncrementId uint32
	dataDir                  string
	idLock                   sync.Mutex
}

func (s *Store) Close() {
	s.Db.Close()
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
	s.Db = db
}

var NeedleMaxAutoIncrementIdKey = []byte("NeedleMaxAutoIncrementIdKey")

func (s *Store) parseNeedleMaxAutoIncrementId() {
	if s.Db == nil {
		logrus.Fatalln("Store's level DB has not initialized yet!")
	}
	maxIdData, err := s.Db.Get(NeedleMaxAutoIncrementIdKey, nil)
	if err != nil {
		s.NeedleMaxAutoIncrementId = 0
		logrus.WithError(err).Infof("store level db has no value for %s", NeedleMaxAutoIncrementIdKey)
	} else {
		s.NeedleMaxAutoIncrementId = binary.LittleEndian.Uint32(maxIdData)
	}
}

func NewStore(ipPort string, dataDir string, stackCount int) *Store {
	b := Store{}
	b.IpPort = ipPort
	b.StackCount = stackCount
	b.StackMaxSize = StackMaxSize
	b.dataDir = dataDir

	b.createLevelDb()
	b.parseNeedleMaxAutoIncrementId()
	for i := 0; i < stackCount; i++ {
		hs := NewHaystack(i, StackMaxSize, b.GetDataDirPath())
		b.Stacks = append(b.Stacks, hs)
	}
	return &b
}

func (s *Store) PutFile(file *multipart.FileHeader) (*Needle, error) {
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
	err = s.Db.Put(needle.IdToByets(), jsonData, nil)
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

func (s *Store) getNeedleFromDb(id []byte) (needle *Needle, err error) {
	jsonData, err := s.Db.Get(id, nil)
	if err != nil {
		return nil, err
	}
	var n = &Needle{}
	proto.Unmarshal(jsonData, n)
	return n, nil
}

func (s *Store) GetFile(id []byte) (n *Needle, err error) {
	n, err = s.getNeedleFromDb(id)
	if err != nil {
		return nil, err
	}
	hs := s.Stacks[n.StackId]
	err = hs.ReadFileBytes(n)
	if err != nil {
		return n, err
	}
	
	return n, nil
}
