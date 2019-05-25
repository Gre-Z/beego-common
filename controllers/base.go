package controllers

import (
	"errors"
	"github.com/Gre-Z/beego-common/show"
	"github.com/Gre-Z/beego-common/validate"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

type BaseControllers struct {
	//继承默认得控制器
	beego.Controller
	//将字符串转成Json格式
	show show.ServerJson
	//参数验证
	//valid *validate.ValiDate
}

func (b *BaseControllers) Prepare() {
	//b.Valid = validate.NewValid(b.Input(), &b.Controller)
}
func (b *BaseControllers) NewValid() *validate.ValiDate {
	return validate.NewValid(b.Input(), &b.Controller)
}

func (b *BaseControllers) Show(code int, msg string, data interface{}, f ...func()) {
	b.show.ServeShow(&b.Controller, code, msg, data, f...)
}

//对请求初步验证符合要求
func requestsInit(b *beego.Controller) (err error) {
	var method = b.Ctx.Input.Method()
	contentType := strings.TrimSpace(b.Ctx.Input.Header("Content-Type"))
	switch method {
	case "POST":
		if contentType != "application/x-www-form-urlencoded" {
			errors.New("content-type error")
			show.ServerJson{}.ServeShow(b, http.StatusForbidden, "content-type error", "")
		}
	case "GET":
	}

	return
}
