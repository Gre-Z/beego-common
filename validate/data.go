package validate

import "github.com/astaxie/beego"

type data struct {
	controller *beego.Controller
}

func (this *data) GetInt(Name string) int {
	i, _ := this.controller.GetInt(Name)
	return i
}

func (this *data) GetInt64(Name string) int64 {
	i, _ := this.controller.GetInt64(Name)
	return i
}

func (this *data) GetBool(Name string) bool {
	b, _ := this.controller.GetBool("bl")
	return b
}

func (this *data) GetString(Name string) string {
	s := this.controller.GetString(Name)
	return s

}

func (this *data) GetFloat(Name string) float64 {
	s, _ := this.controller.GetFloat(Name)
	return s

}
