package validate

import "fmt"

const (
	RgxNumber         = `(-|\+)?\d+`
	RgxTelephone      = "1[345789]\\d{9}"
	RgxAll            = `[\s\S]*`
	RgxBool           = `(true|false)`
	RgxAsc            = `(asc|desc)`
	RgxNumberPositive = `[1-9]\d*$`        //匹配正整数
	RgxChinese        = "[\u4e00-\u9fa5]+" // 匹配中文
)

//这是数字的长度
func RgxSetNumLength(Num uint) string {
	return fmt.Sprintf("(-|\\+)?\\d{%d}", Num)
}

//这是字符串的长度
func RgxSetStrLength(Num uint) string {
	return fmt.Sprintf("[\\s\\S]{%d}", Num)
}
