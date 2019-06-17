package controllers

import (
	"errors"
	"fmt"
	"github.com/Gre-Z/beego-common/show"
	"github.com/Gre-Z/beego-common/validate"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

type BaseControllers struct {
	//继承默认得控制器
	beego.Controller
}

//用于过滤url路径
type _IgnoreLoginer interface {
	IgLogin() []string
}

func (b *BaseControllers) Prepare() {
	b.parseLogin()
}
func (b *BaseControllers) NewValid() *validate.ValiDate {
	return validate.NewValid(b.Input(), &b.Controller)
}

func (b *BaseControllers) Show(code int, msg string, data interface{}, f ...func()) {
	show.NewShow().ServeShow(&b.Controller, code, msg, data, f...)
}

func (b *BaseControllers) parseLogin() {
	token := excludeToken(&b.Controller, b.Ctx.Input.URL())
	if !token {
		//解析token
		fmt.Println("请登录")
		b.StopRun()
	} else {
		//无需解析的呢轮毂
		fmt.Println("无需登陆")
	}
}

func excludeToken(c *beego.Controller, t string) bool {
	var ig []string
	if it, ok := c.AppController.(_IgnoreLoginer); ok {
		ig = it.IgLogin()
	}
	for _, v := range ig {
		isOk := strings.HasSuffix(t, strings.ToLower(v)) || v == "ALL"
		if isOk {
			return isOk
		}
	}
	return false
}

//对请求初步验证符合要求
func requestsInit(b *beego.Controller) (err error) {
	var method = b.Ctx.Input.Method()
	contentType := strings.TrimSpace(b.Ctx.Input.Header("Content-Type"))
	switch method {
	case "POST":
		if contentType != "application/x-www-form-urlencoded" {
			errors.New("content-type error")
			show.NewShow().ServeShow(b, http.StatusForbidden, "content-type error", "")
		}
	case "GET":
	}

	return
}
