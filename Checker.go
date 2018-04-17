package superChecker

import (
	"regexp"
	"fmt"
	"reflect"
	"strings"
	"github.com/pkg/errors"
)

type Checker struct {
	ruler Ruler
}
type Ruler struct {
	RegexBuilder      map[string]*regexp.Regexp
	defaultRegexBuilder map[string]*regexp.Regexp
}

func GetChecker() *Checker{
	checker := &Checker{}
	checker.ruler.defaultRegexBuilder = make(map[string]*regexp.Regexp)
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
		checker.ruler.defaultRegexBuilder[k] = r
	}
	fmt.Println("注入默认成功")
	return checker
}

func (checker *Checker) AddDefaultRegex(key string, regex string) error{
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	key = strings.ToLower(key)
	checker.ruler.defaultRegexBuilder[key] = r
	return nil
}

func (checker *Checker) AddRegex(key string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	key = strings.ToLower(key)
	checker.ruler.RegexBuilder[key] = r
	return nil
}
func (checker *Checker) RemoveRegex(key string) {
	key = strings.ToLower(key)
	delete(checker.ruler.RegexBuilder, key)
}
func (checker *Checker) ListAll() {
	for v, k := range checker.ruler.defaultRegexBuilder {
		fmt.Println(fmt.Sprintf("key:%s,v:%v", v, k))
	}
	for v, k := range checker.ruler.RegexBuilder {
		fmt.Println(fmt.Sprintf("key:%s,v:%v", v, k))
	}
}
func (checker *Checker) ListDefault() {
	for v, k := range checker.ruler.defaultRegexBuilder {
		fmt.Println(fmt.Sprintf("key:%s,v:%v", v, k))
	}
}
func (checker *Checker) ListRegexBuilder() {
	for v, k := range checker.ruler.RegexBuilder {
		fmt.Println(fmt.Sprintf("key:%s,v:%v", v, k))
	}
}
func (checker *Checker) IsContainKey(key string) bool {
	key = strings.ToLower(key)
	for k, _ := range checker.ruler.RegexBuilder {
		if k == key {
			///	fmt.Println("在自定义builder内找到"+key+"匹配规则")
			return true
		}
	}
	for k, _ := range checker.ruler.defaultRegexBuilder {
		if k == key {
			//fmt.Println("在默认builder内找到"+key+"匹配规则")
			return true
		}
	}
	//fmt.Println("没有找到"+key+"匹配规则")
	return false
}

func (checker *Checker) IsBuilderContainKey(key string) bool {
	key = strings.ToLower(key)
	for k, _ := range checker.ruler.RegexBuilder {
		if k == key {
			return true
		}
	}
	return false
}

func (checker *Checker) GetDefaultBuilt() map[string]*regexp.Regexp {
	return checker.ruler.defaultRegexBuilder
}

func (checker *Checker) SuperCheck(input interface{}) (bool, string, error) {
	vType := reflect.TypeOf(input)
	vValue := reflect.ValueOf(input)
	fmt.Println(fmt.Sprintf("input的类型是%v:", vType))
	fmt.Println(fmt.Sprintf("input的值是%v:", vValue))
	for i := 0; i < vType.NumField(); i++ {
		valueStr := vValue.Field(i).String()
		tagValue := vType.Field(i).Tag.Get("superChecker")
		if tagValue==""{
			continue
		}
		tagValue = strings.ToLower(tagValue)
		if strings.Contains(tagValue, "|") {
			if ok, err := rollingCheck(checker, valueStr, tagValue, "|"); !ok {
				if err != nil {
					return false, "检查" + vType.Field(i).Name + "时发生了错误", err
				}
				return false, fmt.Sprintf("%v 匹配失败", vType.Field(i).Name), nil
			}
			//fmt.Println(fmt.Sprintf("%v匹配成功",vType.Field(i).Name))
			continue
		} else {
			if ok, err := rollingCheck(checker, valueStr, tagValue, ","); !ok {
				if err != nil {
					return false, "检查" + vType.Field(i).Name + "时发生了错误", err
				}
				return false, fmt.Sprintf("%v 匹配失败", vType.Field(i).Name), nil
			}
			//fmt.Println(fmt.Sprintf("%v匹配成功",vType.Field(i).Name))

			continue
		}
	}
	return true, "匹配成功", nil
}

func checkRegex(input string, regex *regexp.Regexp) bool {
	return regex.MatchString(input)
}

func rollingCheck(checker *Checker, valueStr string, tagValue string, symbol string) (bool, error) {

	var subStrings = make([]string, 1)
	subStrings = strings.Split(tagValue, symbol)
	for i, v := range subStrings {
		if !checker.IsContainKey(v) {
			return false, errors.New("未定义" + v + "规则")
		}
		if checker.IsBuilderContainKey(v) {
			//fmt.Println("自定义buider包含了"+v+"规则")

			if !checkRegex(valueStr, checker.ruler.RegexBuilder[v]) {
				//fmt.Println(v+"规则匹配失败")
				return false, nil
			} else {
				if symbol == "|" {
					return true, nil
				}
				continue
			}
		}
		if !checkRegex(valueStr, checker.GetDefaultBuilt()[v]) {
			if symbol == "," {
				return false, nil
			} else {
				if i == len(subStrings)-1 {
					return false, nil
				}
				continue
			}
		} else {
			if symbol == "|" {
				return true, nil
			} else {
				if i == len(subStrings)-1 {
					return true, nil
				}
				continue
			}
		}

	}
	return true, nil

}

func(checker *Checker) Check(input string,regex string) (bool,error){
	r, er := regexp.Compile(regex)
	if er!=nil{
		return false,er
	}
	return r.MatchString(input),nil
}