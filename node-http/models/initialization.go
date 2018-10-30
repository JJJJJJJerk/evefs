package models

import (
	"fmt"
	"github.com/dejavuzhou/zerg/config"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

//redis client
var RedisDb *redis.Client

var MysqlDb *gorm.DB

//初始化redis
func init() {
	//initializing redis client
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASSWORD, // no password set
		DB:       config.REDIS_DB_IDX,   // use default DB
	})
	if pong, err := RedisDb.Ping().Result(); err != nil || pong != "PONG" {
		logrus.WithError(err).Error("initializing spider has failled for the result of pinging redis is not PONG or redis configuration is wrong")
	}
	
	conn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", config.MYSQL_USER, config.MYSQL_PASSWORD, config.MYSQL_HOST, config.MYSQL_PORT, config.MYSQL_DB_NAME, config.MYSQL_CHARSET)
	if db, err := gorm.Open("mysql", conn); err == nil {
		//db.LogMode(config.MYSQL_ENABLE_LOG)
		MysqlDb = db
	} else {
		logrus.WithError(err).Error("initialize mysql database failed")
	}
}
