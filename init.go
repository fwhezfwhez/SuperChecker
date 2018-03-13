package superChecker

import (
	"regexp"
	"fmt"
	"strings"
)


var checker *Checker

func init() {
	checker = &Checker{}
	checker.ruler.defaultRegexBuilt = make(map[string]*regexp.Regexp)
	checker.ruler.RegexBuilder =make(map[string]*regexp.Regexp)
	fmt.Println("分配成功")
	regexes := map[string]string{
		"UserName":    "^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$", //中文英文下划线点的组合，长度40以内，是常用的用户名正则限制
		"Number":      "^[0-9]+$",                           //一个以上数字			// 正整数
		"Decimal":     "^\\d+\\.[0-9]+$",                    //小数
		"MobilePhone": "^1[0-9]{10}$",                       //移动电话
		"TelePhone":   "^[0-9]{8}$",                         // 家用电话
		"NotNull":"^[\\s\\S]+$",
	}

	for k,v:=range regexes{
		r, _ := regexp.Compile(v)
		k=strings.ToLower(k)
		checker.ruler.defaultRegexBuilt[k] = r
	}
	fmt.Println("注入默认成功")

}
func GetChecker() *Checker{
	return checker
}