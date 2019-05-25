package validate

import (
	"fmt"
	"github.com/Gre-Z/beego-common/show"
	"github.com/astaxie/beego"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

var mt sync.Mutex
var nonExistent = "【%s】必须存在   /   注释:%s"
var formatErr = "【%s】内容不符合要求   /   注释:%s"

func SetNonExistent(format string) {
	nonExistent = format
}
func SetFormatErr(format string) {
	formatErr = format
}

type ValiDate struct {
	date       map[string]Date
	input      url.Values
	Error      []error
	wg         *sync.WaitGroup
	controller *beego.Controller
	funcList   []func(*sync.WaitGroup)
}

type Date struct {
	Name    string
	Rgx     *regexp.Regexp
	Default interface{}
	Msg     string
}

// valid的工厂函数
func NewValid(values url.Values, controller *beego.Controller) *ValiDate {
	date := ValiDate{
		date:       make(map[string]Date),
		input:      values,
		controller: controller,
		wg:         new(sync.WaitGroup),
	}
	return &date
}

// 添加错误信息
func (this *ValiDate) addErr(err error) {
	this.Error = append(this.Error, err)
}

// 匹配正整数
func (this *ValiDate) Int(Name, Msg string, Default interface{}) *ValiDate {
	this.valid(Name, RgxNumberPositive, Msg, Default)
	return this
}

// 匹配数字范围
func (this *ValiDate) RangeNum(Name string, Min, Max int, Msg string, Default interface{}) *ValiDate {
	rangeNum := func(v string) {
		var err error
		i, err := strconv.Atoi(v)
		if err != nil {
			this.addErr(err)
		}
		if !(Min <= i && i <= Max) {
			this.addErr(fmt.Errorf(formatErr, Name, Msg))
		}
	}
	this.valid(Name, RgxNumber, Msg, Default, rangeNum)
	return this
}

// 通用正则匹配文本
func (this *ValiDate) Normal(Name string, Rgx string, Msg string, Default interface{}) *ValiDate {
	this.valid(Name, Rgx, Msg, Default)
	return this
}

// 字段必须存在
func (this *ValiDate) FieldsMust(Name string, Msg string) *ValiDate {
	this.valid(Name, RgxAll, Msg, nil)
	return this
}

//函数封装 延迟执行
func (this *ValiDate) valid(Name string, Rgx string, Msg string, Default interface{}, op ...func(v string)) {
	this.funcList = append(this.funcList, func(wg *sync.WaitGroup) {
		var err error
		defer func() {
			if err != nil {
				this.Error = append(this.Error, err)
				delete(this.input, Name)
			} else {
				for _, v := range op {
					v(this.input[Name][0])
					//v.Call([]reflect.Value{reflect.ValueOf(this.input[Name][0])})
				}
			}
			defer wg.Done()
		}()
		date := Date{Name: Name, Msg: Msg, Default: Default}
		if Rgx == "" {
			Rgx = RgxAll
		}
		if date.Rgx, err = regexp.Compile("^" + Rgx + "$"); err == nil {
			val, ok := this.input[Name]
			if !ok {
				if date.Default == nil {
					err = fmt.Errorf(nonExistent, date.Name, date.Msg)
				} else {
					valueTemp, bl := date.Default.(string)
					if !bl {
						err = fmt.Errorf("默认值仅支持 (string) 类型")
					} else {
						this.writeMap(date.Name, valueTemp)
					}
				}
			} else {
				if !date.Rgx.MatchString(val[0]) {
					err = fmt.Errorf(formatErr, date.Name, date.Msg)
				}
				//else {
				//	this.input[date.Name] = []string{val[0]}
				//}
			}
		}
	})
}

//协程写入 资源锁
func (this *ValiDate) writeMap(key, value string) {
	mt.Lock()
	this.input[key] = []string{value}
	mt.Unlock()
}

//协程并发 执行函数
func (this *ValiDate) Exec() *ValiDate {
	for _, v := range this.funcList {
		this.wg.Add(1)
		go v(this.wg)
	}
	this.wg.Wait()
	return this
}

//Beego重新设置函数
func (this *ValiDate) ParamReset() {
	this.Exec()
	if len(this.Error) > 0 {
		show.ServerJson{}.ServeShow(this.controller, http.StatusBadRequest, this.Error[0].Error(), "")
	}
	bTemp := this.controller.Ctx.Input
	bTemp.ResetParams()
	for k, v := range this.input {
		bTemp.SetParam(k, v[0])
	}
}

//验证分页
func (this *ValiDate) Pagination() *ValiDate {
	this.valid("page", RgxNumberPositive, "当前页码", nil)
	this.valid("page_size", RgxNumberPositive, "显示数量", nil)
	return this
}
