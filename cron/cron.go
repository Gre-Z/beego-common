package cron

import (
	"errors"
	"fmt"
	"github.com/Gre-Z/common"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	"reflect"
	"strings"
	"sync"
)

// crontab时间遇到999的情况
var cron Cron
var canRegister bool

type TaskMap map[string]TaskDetail
type Cron struct {
	mt      sync.Mutex
	taskMap TaskMap
	debug   bool
}

type CronTask struct {
	common.Model
	TaskName   string `json:"task_name" gorm:"unique_index"`
	Controller string `json:"controller"`
	Func       string `json:"func"`
	Param      string `json:"param"`
	ParamNum   int    `json:"param_num" gorm:"-"`
	Spec       string `json:"spec"`
	Status     int    `gorm:"type:enum('1','2');default:'2'"`
}

// 用来记录该控制器的详情
type TaskDetail struct {
	Func       reflect.Value
	Controller string `json:"controller_name"`
	FuncName   string `json:"func_name"`
	ParamsNum  int    `json:"params_num"`
}

// 初始化map
func init() {
	cron.taskMap = make(TaskMap, 0)
	bl, err := beego.AppConfig.Bool("cron.debug")
	if err != nil {
		cron.debug = false
	} else {
		cron.debug = bl
	}

}

// 将对象注册到任务map
func (this *Cron) RegisterMethod(values ...interface{}) {
	for _, v := range values {
		// 获取对象反射类型
		of := reflect.TypeOf(v).Elem()
		structName := of.Name()
		structPath := of.PkgPath()
		path := structPath + "/" + structName

		// 获取对象反射值
		valueOf := reflect.ValueOf(v)
		// 遍历该对象的方法列表
		for i := 0; i < valueOf.NumMethod(); i++ {
			method := valueOf.Method(i)
			s := valueOf.Type().Method(i).Name
			suffix := strings.HasSuffix(s, "Task")
			if !suffix {
				continue
			}
			taskDetail := TaskDetail{
				Func:       method,
				Controller: path,
				FuncName:   s,
				ParamsNum:  method.Type().NumIn() - 1,
			}
			this.mt.Lock()
			this.taskMap[path+":"+s] = taskDetail
			this.mt.Unlock()
		}
	}
	toolbox.StartTask()
}

// 返回可执行任务列表详情
func (this *Cron) CronMap() TaskMap {
	return this.taskMap
}

// 启动/重启任务
func (this *Cron) Restart(task CronTask) (error) {
	value, e := this.CronIsExist(task.Controller, task.Func)
	if e != nil {
		return e
	}
	fc := func() error {
		if cron.debug {
			this.log(task.TaskName, task)
		}
		var params []reflect.Value
		tasker := toolbox.AdminTaskList[task.TaskName]
		params = append(params, reflect.ValueOf(tasker))
		if task.Param != "" {
			split := strings.Split(task.Param, ",")
			for _, v := range split {
				of := reflect.ValueOf(v)
				params = append(params, of)
			}
		}
		value.Func.Call(params)
		return nil
	}
	go fc()
	newTask := toolbox.NewTask(task.TaskName, task.Spec, fc)
	toolbox.DeleteTask(task.TaskName)
	toolbox.AddTask(task.TaskName, newTask)
	return nil
}

// 停止任务
func (this *Cron) Stop(task CronTask) {
	toolbox.DeleteTask(task.TaskName)
}

// 判断控制器函数是否存在 存在则返回详情
func (this *Cron) CronIsExist(path, name string) (TaskDetail, error) {
	if v, ok := this.taskMap[path+":"+name]; ok {
		return v, nil
	}
	return TaskDetail{}, errors.New("cron is not exist")
}

// 判断任务是否已经存在
func (this *Cron) TaskIsExist(taskName string) (toolbox.Tasker, bool) {
	if v, ok := toolbox.AdminTaskList[taskName]; ok {
		return v, ok
	} else {
		return &toolbox.Task{}, ok
	}
}

// 输出当前任务的详情
func (this *Cron) log(taskName string, task CronTask) {
	tasker, bl := this.TaskIsExist(taskName)
	if !bl {
		logs.Error("任务不存在")
	}
	taskObj := InterfaceToObject(tasker)
	logs.Warning(fmt.Sprintf("\n任务名:\t%s\n执行路径:%s\n函数名:%s\n当前执行时间:%s\n下次执行时间:%s\n=====", taskObj.Taskname, task.Controller, task.Func, taskObj.Prev, taskObj.Next))
}

// 接口转实例化对象
func InterfaceToObject(tasker toolbox.Tasker) *toolbox.Task {
	task := tasker.(*toolbox.Task)
	return task
}

// 工厂函数
func CronObj() *Cron {
	return &cron
}
