package controllers

import (
	"errors"
	"github.com/Gre-Z/beego-common/show"
	"github.com/Gre-Z/beego-common/validate"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
	"strings"
)

type BaseControllers struct {
	//继承默认得控制器
	beego.Controller
}

// 用于过滤url路径
type _Ignore interface {
	Ignore() []string
}

// 用户 返回 token
type _JwtParser interface {
	JwtParse(string) (error)
}

func (b *BaseControllers) Prepare() {
	b.parseLogin()
}

func (b *BaseControllers) NewValid() *validate.ValiDate {
	return validate.NewValid(b.Input(), &b.Controller)
}

func (b *BaseControllers) Show(code int, msg interface{}, data interface{}, f ...func()) {
	show.NewShow().ServeShow(&b.Controller, code, msg, data, f...)
}

func (b *BaseControllers) parseLogin() {
	jump := excludeLogin(&b.Controller, b.Ctx.Input.URL())
	if !jump {
		if it, ok := b.AppController.(_JwtParser); ok {
			Authorization := b.Ctx.Request.Header.Get("Authorization")
			err := it.JwtParse(Authorization)
			if err != nil {
				b.Show(http.StatusForbidden, err, nil)
			}
		}
	}
}

func excludeLogin(c *beego.Controller, t string) bool {
	var ig []string
	if it, ok := c.AppController.(_Ignore); ok {
		ig = it.Ignore()
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

func AccessControllAllowOrigin() func(ctx *context.Context) bool {
	return func(ctx *context.Context) bool {
		origin := ctx.Request.Header.Get("Origin")
		allowOrigin := []string{}
		allow := beego.AppConfig.String("Access-Control-Allow-Methods")
		if allow != "" {
			allowOrigin = strings.Split(allow, ",")
		}
		ctx.Output.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS")                               //允许post访问
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-XSRF-TOKEN") //header的类型
		ctx.Output.Header("Access-Control-Max-Age", "1728000")
		ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		res := false
		for _, e := range allowOrigin {
			if strings.EqualFold(e, origin) {
				res = true
				break
			}
		}
		if res {
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
		}
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(204)
			ctx.Output.Body([]byte(""))
			return false
		}
		return true
	}
}
