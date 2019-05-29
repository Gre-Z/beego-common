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
	user := beego.AppConfig.String("mysql.user")
	password := beego.AppConfig.String("mysql.password")
	addr := beego.AppConfig.String("mysql.host")
	dbname := beego.AppConfig.String("mysql.dbname")
	singularTable := beego.AppConfig.DefaultBool("mysql.singularTable", false)
	logMode := beego.AppConfig.DefaultBool("mysql.logMode", false)
	maxIdle := beego.AppConfig.DefaultInt("mysql.maxIdle", 0)
	maxOpen := beego.AppConfig.DefaultInt("maxOpen", 0)

	config.Mysql = mysql.Options{
		User: user, Password: password, Addr: addr, Dbname: dbname,
		SingularTable: singularTable,
		LogMode:       logMode,
		MaxIdle:       maxIdle,
		MaxOpen:       maxOpen,
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
