package show

import (
	"github.com/astaxie/beego"
	"github.com/jessonchan/jsun"
)

type ServerJson struct{}

func NewShow() *ServerJson {
	return new(ServerJson)
}

func (s *ServerJson) ServeShow(c *beego.Controller, code int, msg interface{}, data interface{}, f ...func()) {
	output := c.Ctx.Output
	retData := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	jsun.UnderScoreStyle()
	content, _ := jsun.Marshal(&retData)
	output.Header("Content-Type", "application/json; charset=utf-8")
	output.Body(content)
	//c.Data["json"] = retData
	//c.ServeJSON()
	for i := range f { //预留日志等操作
		go f[i]()
	}
	c.StopRun()
}
