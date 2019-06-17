package validate

import "github.com/astaxie/beego"

var nonExistent = "字段【%s】:必须存在   /   注释:%s"
var formatErr = "字段【%s】:值不符合   /   注释:%s"
var isDebug = true

func SetNonExistent(format string) {
	nonExistent = format
}
func SetFormatErr(format string) {
	formatErr = format
}
func SetIsDebug(bl bool) {
	isDebug = bl
}

func init() {
	isDebug = beego.AppConfig.DefaultBool("validate.isDebug", false)
}
