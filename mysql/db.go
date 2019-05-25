/*
mysql.user = 用户名
mysql.password = 密码
mysql.host = ip地址
mysql.port = 端口
mysql.dbname = 数据库名字
gorm.singularTable = 全局禁用复数
gorm.logMode = 开启日志
gorm.maxOpen = 最大打开的连接数 0表示不限制
gorm.maxIdle = 闲置的连接数量
*/
package mysql

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//gorm model
type Model struct {
	Id        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
type DB struct {
	Default *gorm.DB
}

var eor error
var db DB

func init() {
	db.Default, eor = gorm.Open("mysql",
		beego.AppConfig.String("mysql.user")+
			":"+
			beego.AppConfig.String("mysql.password")+
			"@tcp("+
			beego.AppConfig.String("mysql.host")+
			")/"+
			beego.AppConfig.String("mysql.dbname")+
			"?charset=utf8mb4&parseTime=true&loc=Local")
	if eor != nil {
		panic(eor)
		return
	} else {
		logs.Info("mysql connect success")
	}

	if singularTable, err := beego.AppConfig.Bool("gorm.singularTable"); err == nil {
		db.Default.SingularTable(singularTable)
	}

	if logMode, err := beego.AppConfig.Bool("gorm.logMode"); err == nil {
		db.Default.LogMode(logMode)
	}

	if maxIdle, err := beego.AppConfig.Int("gorm.maxIdle"); err == nil {
		db.Default.DB().SetMaxIdleConns(maxIdle)
	}

	if maxOpen, err := beego.AppConfig.Int("gorm.maxOpen"); err == nil {
		db.Default.DB().SetMaxOpenConns(maxOpen)
	}
}

func (DB) MysqlNew() *gorm.DB {
	if db.Default == nil {
		logs.Info("连接错误")
	}
	return db.Default.New()
}
