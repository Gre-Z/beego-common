package redis

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

type RD struct {
	myDefault *redis.Client
}

var rd RD
var Eor error

func Init() {
	dbname, Eor := beego.AppConfig.Int("redis.dbname")
	if Eor != nil {
		dbname = 0
	}
	NewClient(dbname)
}

func NewClient(dbname int) {
	rd.myDefault = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.String("redis.addr"),
		Password: beego.AppConfig.String("redis.password"),
		DB:       dbname,
	})
	_, Eor = rd.myDefault.Ping().Result()
	if Eor != nil {
		panic(Eor)
		return
	} else {
		logs.Info("redis connect success")
	}
}

func (RD) RedisNew() *redis.Client {
	return rd.myDefault
}

func ChangeClient(client *redis.Client) RD {
	rd.myDefault = client
	return rd
}
