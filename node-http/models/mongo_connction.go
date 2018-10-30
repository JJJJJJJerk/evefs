package models

import "fmt"

const (
	SshAuthTypePrivateKey = iota
	SshAuthTypePassword
)

type SshTunnelModel struct {
	IsUseSshTunnel          bool
	SshAddress              string
	SshPort                 string
	SshUserName             string
	SshAuthType             int
	SshPrivateKey           string
	SshPrivateKeyPassphrase string
}

type MongoConnectionModel struct {
	ConnectionName string
	
	Address       string
	Port          string
	UserName      string
	Password      string
	DefaultDbName string
	
	SshTunnelModel
}

func (m *MongoConnectionModel) ConnectionString() string {
	//https://docs.mongodb.com/manual/reference/connection-string/index.html#standard-connection-string-format
	//mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	return fmt.Sprintf(`mongodb://%s:%s@%s:%s/%s`, m.UserName, m.Password, m.Address, m.Port, m.DefaultDbName)
}
