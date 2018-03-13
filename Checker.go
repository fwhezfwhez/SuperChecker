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
	defaultRegexBuilt map[string]*regexp.Regexp
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
	for v, k := range checker.ruler.RegexBuilder {
		fmt.Sprintln("key:%s,v:%v", v, k)
	}
}


func (checker *Checker) IsContainKey(key string) bool {
	key = strings.ToLower(key)
	for k, _ := range checker.ruler.RegexBuilder {
		if k == key {
			fmt.Println("在自定义builder内找到"+key+"匹配规则")
			return true
		}
	}
	for k, _ := range checker.ruler.defaultRegexBuilt {
		if k == key {
			fmt.Println("在默认builder内找到"+key+"匹配规则")
			return true
		}
	}
	fmt.Println("没有找到"+key+"匹配规则")
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
	return checker.ruler.defaultRegexBuilt
}

func (checker *Checker) SuperCheck(input interface{}) (bool, error) {
	vType := reflect.TypeOf(input)
	vValue := reflect.ValueOf(input)
	fmt.Println(fmt.Sprintf("input的类型是%v:", vType))
	for i := 0; i < vType.NumField(); i++ {
		valueStr := vValue.Field(i).String()
		tagValue := vType.Field(i).Tag.Get("superChecker")
		tagValue = strings.ToLower(tagValue)
		if strings.Contains(tagValue, ",") {
			if ok, err := rollingCheck(checker, valueStr, tagValue, ","); !ok {
				if err != nil {
					return false, err
				}
				return false, errors.New(fmt.Sprintf("%v 匹配失败", vType.Field(i).Name))
			}
			fmt.Println(fmt.Sprintf("%v匹配成功",vType.Field(i).Name))
			continue
		} else {
			if ok, err := rollingCheck(checker, valueStr, tagValue, "|"); !ok {
				if err != nil {
					return false, err
				}
				return false, errors.New(fmt.Sprintf("%v 匹配失败", vType.Field(i).Name))
			}
			fmt.Println(fmt.Sprintf("%v匹配成功",vType.Field(i).Name))
			continue
		}
	}
	return true, nil
}

func checkRegex(input string, regex *regexp.Regexp) bool {
	return regex.MatchString(input)
}

func rollingCheck(checker *Checker, valueStr string, tagValue string, symbol string) (bool, error) {

	var subStrings = make([]string, 1)
	subStrings = strings.Split(tagValue, symbol)
	for _, v := range subStrings {
		if !checker.IsContainKey(v) {
			return false, errors.New("未定义" + v + "规则")
		}
		if checker.IsBuilderContainKey(v) {
			fmt.Println("自定义buider包含了"+v+"规则")

			if !checkRegex(valueStr, checker.ruler.RegexBuilder[v]) {
				fmt.Println(v+"规则匹配失败")
				return false, nil
			} else {
				if symbol == "|" {
					return true, nil
				}
				continue
			}
		}
		if !checkRegex(valueStr, checker.GetDefaultBuilt()[v]) {
			if symbol=="," {
				return false, nil
			}
			continue
		} else {
			if symbol == "|" {
				return true, nil
			}
		}

	}
	return true, nil

}
