package config

import (
	"github.com/Gre-Z/common/mysql"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

type Config struct {
	Redis *redis.Options
	Mysql mysql.Options
}

func Init() (config Config) {
	//初始化redis
	initRedis(&config)

	//初始化mysql
	initMysql(&config)
	return
}

func initMysql(config *Config) {
	Add := func(name string) string {
		return "mysql." + name
	}
	user := beego.AppConfig.String(Add("user"))
	password := beego.AppConfig.String(Add("password"))
	addr := beego.AppConfig.String(Add("addr"))
	dbname := beego.AppConfig.String(Add("dbname"))
	singularTable := beego.AppConfig.DefaultBool(Add("singularTable"), false)
	logMode := beego.AppConfig.DefaultBool(Add("logMode"), false)
	maxIdle := beego.AppConfig.DefaultInt(Add("maxIdle"), 0)
	maxOpen := beego.AppConfig.DefaultInt(Add("maxOpen"), 0)
	connectName := beego.AppConfig.DefaultString("connectName", "default")
	autoMigrate := beego.AppConfig.DefaultBool("autoMigrate", false)
	config.Mysql = mysql.Options{
		User: user, Password: password, Addr: addr, Dbname: dbname,
		SingularTable: singularTable,
		LogMode:       logMode,
		MaxIdle:       maxIdle,
		MaxOpen:       maxOpen,
		ConnectName:   connectName,
		AutoMigrate:   autoMigrate,
	}

}

func initRedis(config *Config) {
	addr := beego.AppConfig.String("redis.addr")
	password := beego.AppConfig.String("redis.password")
	db, err := beego.AppConfig.Int("redis.dbname")
	if err != nil {
		db = 0
	}
	config.Redis = &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}
