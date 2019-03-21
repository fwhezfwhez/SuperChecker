package superChecker

import (
	"errors"
	"fmt"
	"github.com/fwhezfwhez/jsoncrack"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DEBUG   = 1
	RELEASE = 2
)

// global checker object
// a checker contains its rule which stores the regex rule pool of the default pool and the added pool
type Checker struct {
	l     sync.RWMutex
	mode  int
	ruler Ruler
}

// rule object contained in a checker,it consists of a default pool and an added pool
type Ruler struct {
	defaultLock         sync.RWMutex
	addedLock           sync.RWMutex
	RegexBuilder        map[string]*regexp.Regexp
	defaultRegexBuilder map[string]*regexp.Regexp
	Funcs               map[string]Func
}

// string array that stand for int type
var intTypes = []string{"int", "int16", "int32", "int64", "int8"}

// string array that stand for floatTypes
var floatTypes = []string{"float", "float32", "float64", "decimal"}

// flag range
var flagRange = []string{"range", "in", "regex", "int", "int32", "int64", "int8", "string", "char", "float", "float32", "float64", "decimal", "time.time", "func", "function"}

// Func has a value and its desgin path.
//     value serves for a self design function that deals with the input data 'in interface{}', and returns its result 'ok bool',
//     'message string', 'e error'.
//     path serves for logging where the function is design
// for example:
//     value:
//     func ValideMoney(in interface{}) (bool,string,error){
// 	        v, ok :=in.(float64)
//          if !ok{
//              return false, fmt.Sprintf( want float64 type, got '%v'", in), errors.New(fmt.Sprintf(" want float64 type, got '%v'", in))
//          }
//	    	return true,"success",nil
//     }
//    path:
//    xxx/xxx/xx/main.go: 90
type Func struct {
	Value func(in interface{}, filedName string) (bool, string, error)
	Path  string
}

// get a checker object which has contained regex rule below:
// username : ^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$
// number : "^[0-9]+$"
// decimal : "^\\d+\\.[0-9]+$"
// mobile phone : "^1[0-9]{10}$"
// telephone : "^[0-9]{8}$"
// notnull: "^[\\s\\S]+$"
func GetChecker() *Checker {
	checker := &Checker{}
	checker.ruler.defaultRegexBuilder = make(map[string]*regexp.Regexp)
	checker.ruler.RegexBuilder = make(map[string]*regexp.Regexp)
	checker.ruler.Funcs = make(map[string]Func)

	checker.ruler.defaultRegexBuilder = compiledMap
	return checker
}

// set its mode of superChecker.DEBUG,superChecker.RELEASE
// DEBUG =1
// RELEASE = 2
func (checker *Checker) SetMode(mode int) {
	checker.l.Lock()
	defer checker.l.Unlock()
	checker.mode = mode
}

// add default regex rule into default pool , when the key is already existed, then it will be replaced by the new one
func (checker *Checker) AddDefaultRegex(key string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	key = strings.ToLower(key)
	checker.ruler.defaultLock.Lock()
	defer checker.ruler.defaultLock.Unlock()

	checker.ruler.defaultRegexBuilder[key] = r
	return nil
}

// add regex into added pool.
func (checker *Checker) AddRegex(key string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	key = strings.ToLower(key)

	checker.ruler.addedLock.Lock()
	defer checker.ruler.addedLock.Unlock()
	checker.ruler.RegexBuilder[key] = r
	return nil
}

// remove regex kv from the  added pool.
func (checker *Checker) RemoveRegex(key string) {
	key = strings.ToLower(key)

	checker.ruler.addedLock.Lock()
	defer checker.ruler.addedLock.Unlock()

	delete(checker.ruler.RegexBuilder, key)
}

// get a regex string format, added pool has higher privilege
func (checker *Checker) GetRule(key string) string {
	checker.ruler.addedLock.RLock()
	defer checker.ruler.addedLock.RUnlock()
	v1, ok1 := checker.ruler.RegexBuilder[key]
	if ok1 {
		return strconv.QuoteToASCII(v1.String())
	}
	checker.ruler.defaultLock.RLock()
	defer checker.ruler.defaultLock.RUnlock()
	v2, ok2 := checker.ruler.defaultRegexBuilder[key]
	if ok2 {
		return strconv.QuoteToASCII(v2.String())
	}
	return ""
}

// list all regex compiled in both the default and the added pool.
func (checker *Checker) ListAll() {
	fmt.Println(fmt.Sprintf(" key | value "))

	for k, v := range checker.ruler.defaultRegexBuilder {
		fmt.Println(fmt.Sprintf(` %s | %s `, k, strconv.QuoteToASCII(v.String())))
	}
	for k, v := range checker.ruler.RegexBuilder {
		fmt.Println(fmt.Sprintf(` %s | %s `, k, strconv.QuoteToASCII(v.String())))
	}

	for k, v := range checker.ruler.Funcs {
		fmt.Println(fmt.Sprintf(` %s | %s `, k, strconv.QuoteToASCII(v.Path)))
	}
}

// list default pool
func (checker *Checker) ListDefault() {
	fmt.Println(fmt.Sprintf(" key | value "))

	for k, v := range checker.ruler.defaultRegexBuilder {
		fmt.Println(fmt.Sprintf(" %s | %s ", k, strconv.QuoteToASCII(v.String())))
	}
}

// list added pool
func (checker *Checker) ListRegexBuilder() {
	fmt.Println(fmt.Sprintf(" key | value "))

	for k, v := range checker.ruler.RegexBuilder {
		fmt.Println(fmt.Sprintf(" %s | %s ", k, strconv.QuoteToASCII(v.String())))
	}
}

// whether the key is contained
func (checker *Checker) IsContainKey(key string) bool {
	key = strings.ToLower(key)
	_, ok1 := checker.ruler.RegexBuilder[key]
	_, ok2 := checker.ruler.defaultRegexBuilder[key]
	if ok1 || ok2 {
		return true
	}
	return false
}

// whether the added pool contains the rule key
func (checker *Checker) IsBuilderContainKey(key string) bool {
	key = strings.ToLower(key)

	_, ok := checker.ruler.RegexBuilder[key]
	if ok {
		return true
	}
	return false
}

// whether the func pool contains the func key
func (checker *Checker) ContainFunc(key string) bool {
	key = strings.ToLower(key)
	_, ok := checker.ruler.Funcs[key]
	return ok
}

// get func by key
func (checker *Checker) GetFunc(key string) Func {
	key = strings.ToLower(key)
	return checker.ruler.Funcs[key]
}

// add a func into func pool
// keyAndPath stands for the func's key and func's define path.
// key must specific and must be keyAndPath[0], path is optional.
// when the length of keyAndPath is 0 or >2 , then throws error.
// when the length of keyAndPath is 1, key is keyAndPath[0], path is the caller stack depth 1.
// when the length of keyAndPath is 2, key is keyAndPath[0], path is keyAndPath[1].
func (checker *Checker) AddFunc(f func(in interface{}, fieldName string) (bool, string, error), keyAndPath ...string) error {
	if len(keyAndPath) > 2 {
		return errors.New(fmt.Sprintf("keyAndPath should has length no more than 2, but got %v", keyAndPath))
	} else if len(keyAndPath) == 0 {
		return errors.New("keyAndPath should at least has length by 1 to define its key, but got 0")

	}
	key := keyAndPath[0]
	var path = ""
	if len(keyAndPath) == 2 {
		path = keyAndPath[1]
	} else if len(keyAndPath) == 1 {
		_, file, line, _ := runtime.Caller(1)
		path = fmt.Sprintf("%s:%d", file, line)
	}
	checker.ruler.Funcs[strings.ToLower(key)] = Func{
		Value: f,
		Path:  path,
	}
	return nil
}

// get the default pool
func (checker *Checker) GetDefaultBuilder() map[string]*regexp.Regexp {
	checker.ruler.defaultLock.RLock()
	defer checker.ruler.defaultLock.RUnlock()
	return checker.ruler.defaultRegexBuilder
}

// Latest version has stop this's developing replaced by functions like 'FormatCheck' 'Validate' 'ValidateMethods' which alias with tag 'validate'
// However the old usages are remained, function is access if you insist on using it
//
// support for string input or type that can be transfer to a string or an object which has function String().
// notice:
// 1. the value of tag 'superCheck' can be either upper or lower or mixed,
//    `superChecker:"username"`,`superChecker:"usERName"`  are ok
// 2. some cases will be ignored:
//    `superChecker:""`, `superChecker:"-"` will be ignored
//     struct{name string}{name:"undefined"}, struct{name string}{name:"undefine"} will be ignored
// 3. make sure the not-ignored fields is string-able, these types can be well stringed:
//    [int,int8,int16,int32,int64,float32,float64,] || <object'function String()>
func (checker *Checker) SuperCheck(input interface{}) (bool, string, error) {
	vType := reflect.TypeOf(input)
	vValue := reflect.ValueOf(input)

	if checker.mode == DEBUG {
		SmartPrint(input)
	}
	var valueStr = ""
	for i := 0; i < vValue.NumField(); i++ {
		tagValue := vType.Field(i).Tag.Get("superChecker")
		//`superChecker:"username"`,`superChecker:"usERName"`  are ok
		tagValue = strings.ToLower(tagValue)
		//`superChecker:""`, `superChecker:"-"` will be ignored
		if tagValue == "" || tagValue == "-" {
			continue
		}

		value := vValue.Field(i).Interface()

		valueStr = ToString(value)

		if valueStr == "undefined" || valueStr == "undefine" {
			continue
		}

		// when contains '|'
		if strings.Contains(tagValue, "|") {
			if ok, err := rollingCheck(checker, valueStr, tagValue, "|"); !ok {
				if err != nil {
					return false, fmt.Sprintf("checking '%s' catch an error '%s'", vType.Field(i).Name, err.Error()), err
				}

				return false, fmt.Sprintf("'%s' unmatched, expected rule '%s',got '%s'", vType.Field(i).Name, checker.GetRule(tagValue), value), nil
			}
			//fmt.Println(fmt.Sprintf("field '%s' success",vType.Field(i).Name))
			//continue
		} else {
			// when contains ',' or neither contains ',' or '|'
			if ok, err := rollingCheck(checker, valueStr, tagValue, ","); !ok {
				if err != nil {
					return false, fmt.Sprintf("checking '%s' catch an error '%s'", vType.Field(i).Name, err.Error()), err
				}
				return false, fmt.Sprintf("'%s' unmatched, expected rule '%s',got '%s'", vType.Field(i).Name, checker.GetRule(tagValue), value), nil
			}
			//fmt.Println(fmt.Sprintf("field '%s' success",vType.Field(i).Name))

			//continue
		}
		continue
	}
	return true, "success", nil
}

// validate if an input value is correct or not
// notice:
//     1. some ignored cases:
//          `validate:""`, `validate:"-"` will be ignored
//          struct{name string}{name:"undefine"}, struct{name string}{name:"undefined"} will be ignored
// support int types ,float types, string, time
func (checker *Checker) FormatCheck(input interface{}) (bool, string, error) {
	vType := reflect.TypeOf(input)
	vValue := reflect.ValueOf(input)
	valueStr := ""
	var ok, ok1, ok2, ok3 bool
	if checker.mode == DEBUG {
		SmartPrint(input)
	}
L:
	for i := 0; i < vType.NumField(); i++ {
		tagValue := vType.Field(i).Tag.Get("validate")
		if whetherLowerCase(tagValue) {
			tagValue = strings.ToLower(tagValue)
		}

		// `validate:""` `validate:"-"` will be ignored
		if tagValue == "" || tagValue == "-" {
			continue
		}
		value := vValue.Field(i).Interface()

		// the empty value will be ignore if no 'notnull' flag
		if !strings.Contains(strings.ToLower(tagValue), "notnull") {
			var zeroTime time.Time
			if ToString(value) == "" || ToString(value) == zeroTime.Format("2006-01-02 15:04:05") {
				continue
			}
		}

		tmp := strings.Split(tagValue, ",")
		flag := strings.ToLower(tmp[0])
		if !strings.Contains(strings.Join(flagRange, " "), flag) {
			return false,
				fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", vType.Field(i).Name, flag, flagRange),
				errors.New(fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", vType.Field(i).Name, flag, flagRange),
				)
		}

		var rule string
		if len(tmp) > 1 {
			rule = strings.Join(tmp[1:], ",")
		}

		// flag range/in validate
		if flag == "range" || flag == "in" {
			if !strings.HasPrefix(rule, "[") || !strings.HasSuffix(rule, "]") {
				return false,
					fmt.Sprintf("field '%s' range/in flag must have its rule format like '[x,x,x,x,x] but got '%s'", vType.Field(i).Name, rule),
					errors.New(fmt.Sprintf("field '%s' range/in flag must have its rule format like '[x,x,x,x,x] but got '%s'", vType.Field(i).Name, rule))
			}
			valueStr = ToString(value)
			arr := strings.Split(rule[1:len(rule)-1], ",")
			for _, v := range arr {
				if v == valueStr {
					continue L
				}
			}
			return false, fmt.Sprintf("field '%s' required in '%s' but got '%s'", vType.Field(i).Name, rule, valueStr), nil
		}

		// regex validate
		// regex validate is used to replace superChecker tag
		// Name string `superChecker:"key"`  <==> Name string `validate:"regex,key"`
		if flag == "regex" {
			if whetherLowerCase(tagValue) {
				rule = strings.ToLower(rule)
			}
			valueStr = ToString(value)
			// rule is raw regex
			if strings.HasPrefix(rule, `^`) && strings.HasSuffix(rule, `$`) {
				ok, er := checker.Check(valueStr, rule)
				if er != nil {
					return false,
						fmt.Sprintf("while validating field '%s' regex '%s' throws an error '%s'", vType.Field(i).Name, strconv.QuoteToASCII(rule), er.Error()),
						errors.New(fmt.Sprintf("while validating field '%s' regex '%s' throws an error '%s'", vType.Field(i).Name, strconv.QuoteToASCII(rule), er.Error()))
				}
				if !ok {
					return false,
						fmt.Sprintf("while validating field '%s' regex '%s' but got unmatched value '%s'", vType.Field(i).Name, strconv.QuoteToASCII(rule), valueStr),
						nil
				}
				continue L
			} else {
				if strings.Contains(rule, "|") {
					// rule formated like 'key1|key2|key3' which can be separated by '|'
					rules := strings.Split(rule, "|")
					for j, v := range rules {
						ok, er := checker.CheckFromPool(valueStr, v)
						if er != nil {
							return false,
								fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", vType.Field(i).Name, j, v, er.Error()),
								errors.New(fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", vType.Field(i).Name, j, v, er.Error()))
						}
						if ok {
							continue L
						}
						if j == len(rules)-1 {
							return false, fmt.Sprintf("while validating field '%s', regex key group %v all fail, got unmatched value '%s'", vType.Field(i).Name, rules, valueStr), nil
						}
					}
				} else if strings.Contains(rule, ",") {
					// rule formated like 'key1,key2,key3' which can be separated by ','
					rules := strings.Split(rule, ",")
					for i, v := range rules {
						ok, er := checker.CheckFromPool(valueStr, v)
						if er != nil {
							return false,
								fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", vType.Field(i).Name, i, v, er.Error()),
								errors.New(fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", vType.Field(i).Name, i, v, er.Error()))
						}
						if !ok {
							return false,
								fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' but got unmatched value '%s'", vType.Field(i).Name, i, v, valueStr),
								nil
						}
					}
				} else {
					// rule is regarded as the key itself
					if !checker.IsContainKey(rule) {
						return false,
							fmt.Sprintf("while validating field '%s' regex key '%s' not found in any of default regex pool or add regex pool,you may use '%s' before using it", vType.Field(i).Name, rule, "checker.AddRegex('key', '^raw regex$')"),
							errors.New(fmt.Sprintf("while validating field '%s' regex key '%s' not found in any of default regex pool or add regex pool,you may use '%s' before using it", vType.Field(i).Name, rule, "checker.AddRegex('key', '^raw regex$')"))
					}
					ok, er := checker.CheckFromPool(valueStr, rule)
					if er != nil {
						return false,
							fmt.Sprintf("while validating field '%s' regex pool key '%s' throws an error '%s'", vType.Field(i).Name, rule, er.Error()),
							er
					}
					if !ok {
						return false,
							fmt.Sprintf("while validating field '%s' regex pool key '%s' but got unmatched value '%s'", vType.Field(i).Name, rule, valueStr),
							nil
					}
				}
			}
			continue L
		}

		// type validate including:
		// func validate,
		// int,float,time validate
		_, ok1 = value.(time.Time)
		_, ok3 = value.(jsoncrack.Time)
		ok2 = in(strings.Split(tagValue, ",")[0],"time.time", "time.Time")
		ok = ok1 || ok2 || ok3
		if ok && strings.Contains(tagValue, ",") && strings.Split(tagValue, ",")[1] != "" {
			valueStr = ToString(value, strings.Split(tagValue, ",")[1])
		} else {
			valueStr = ToString(value)
		}

		if valueStr == "undefined" || valueStr == "undefine" || valueStr == "" {
			continue
		}

		if strings.Contains(tagValue, ",") {
			tmp = strings.Split(tagValue, ",")
			tagValue = tmp[0]
			rule = strings.Join(tmp[1:], ",")
			if isFunc(tagValue) {
				if len(tmp) < 2 {
					return false,
						fmt.Sprintf("'%s' is validated as 'func', the tag 'validate' must has its tag value length more than 2,but got '%s' length is %d", vType.Field(i).Name, tagValue, len(tmp)),
						errors.New(fmt.Sprintf("'%s' is validated as 'func', the tag 'validate' must has its tag value length more than 2,but got '%s' length %d", vType.Field(i).Name, tagValue, len(tmp)))
				}

				if len(tmp) == 2 {
					if arr := strings.Split(rule, "|"); len(arr) > 1 {
						// validate:func,key1|key2|key3
						for j, r := range arr {
							if !checker.ContainFunc(r) {
								return false, fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", vType.Field(i).Name, j, r),
									errors.New(fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", vType.Field(i).Name, j, r))
							}
							ok, msg, er := checker.GetFunc(r).Value(value, vType.Field(i).Name)
							if ok {
								continue L
							}
							if j >= len(arr)-1 {
								return ok, msg, er
							}
							continue
						}
					} else {
						// validate:func,key1
						if !checker.ContainFunc(rule) {
							return false, fmt.Sprintf("'%s' func has not be added into func pool,use checker.AddFunc() to register", rule),
								errors.New(fmt.Sprintf("'%s' func has not be added into func pool,use checker.AddFunc() to register", rule))
						}
						ok, msg, er := checker.GetFunc(rule).Value(value, vType.Field(i).Name)
						if ok {
							continue L
						}
						return ok, msg, er
					}
				} else {
					// validate:func,key1,key2,key3,key4
					rules := strings.Split(rule, ",")
					for j, r := range rules {
						if !checker.ContainFunc(r) {
							return false, fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", vType.Field(i).Name, j, r),
								errors.New(fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", vType.Field(i).Name, j, r))
						}
						ok, msg, er := checker.GetFunc(r).Value(value, vType.Field(i).Name)
						if !ok {
							return ok, msg, er
						}
						if j == len(rules)-1 {
							continue L
						}
					}
				}

			} else if IsInt(tagValue) {
				if rule != "" {
					tmp2 := strings.Split(rule, ":")
					if len(tmp2) != 2 {
						return false, "", errors.New("notation requires number1:number2,but got " + rule)
					}
					v, er := strconv.Atoi(valueStr)
					if er != nil {
						return false, vType.Field(i).Name + " format required int but got " + valueStr, nil
					}
					if tmp2[0] != "" {
						min, er := strconv.Atoi(tmp2[0])
						if er != nil {
							return false, "", errors.New(vType.Field(i).Name + " notation rule required number:number but get " + tmp2[0])
						}
						if v < min {
							return false, vType.Field(i).Name + " int value required bigger than " + tmp2[0] + " but get " + valueStr, nil
						}
					}
					if tmp2[1] != "" {
						max, er := strconv.Atoi(tmp2[1])
						if er != nil {
							return false, "", errors.New(vType.Field(i).Name + " notation rule required number:number but get " + tmp2[1])
						}
						if v > max {
							return false, vType.Field(i).Name + " int value required smaller than " + tmp2[1] + " but get " + valueStr, nil
						}
					}

				} else {
					_, er := strconv.Atoi(valueStr)
					if er != nil {
						return false, vType.Field(i).Name + " format required int but got " + valueStr, nil
					}
				}
			} else if IsFloat(tagValue) {
				if rule != "" {
					tmp2 := strings.Split(rule, ":")
					if len(tmp2) != 2 {
						return false, "", errors.New(" notation requires float_number1:float_number2,but got " + rule)
					}
					v, er := strconv.ParseFloat(valueStr, 64)
					if er != nil {
						return false, vType.Field(i).Name + " format required float but got " + valueStr, nil
					}
					if tmp2[0] != "" {
						min, er := strconv.ParseFloat(tmp2[0], 64)
						if er != nil {
							return false, "", errors.New(vType.Field(i).Name + " notation rule required float_number:float_number but got " + tmp2[0])
						}
						if v < min {
							return false, vType.Field(i).Name + " float value required bigger than" + tmp2[0] + " but got " + valueStr, nil
						}
					}
					if tmp2[1] != "" {
						max, er := strconv.ParseFloat(tmp2[1], 64)
						if er != nil {
							return false, "", errors.New(vType.Field(i).Name + " notation rule required number:number but got " + tmp2[1])
						}
						if v > max {
							return false, vType.Field(i).Name + " float value required smaller than " + tmp2[1] + " but got " + valueStr, nil
						}
					}
				} else {
					_, er := strconv.ParseFloat(valueStr, 64)
					if er != nil {
						return false, vType.Field(i).Name + "format required float but got " + valueStr, nil
					}
				}
			} else if in(tagValue , "time.time", "time.Time") {
				//"2006/1/2 15:04:05"
				if rule != "" {
					_, er := time.ParseInLocation(rule, valueStr, time.Local)
					if er != nil {
						return false, fmt.Sprintf("while validating field '%s', time format requires %s but go %s", vType.Field(i).Name, rule, valueStr), nil
					}
				} else {
					_, er := time.ParseInLocation("2006/1/2 15:04:05", valueStr, time.Local)
					if er != nil {
						return false, fmt.Sprintf("while validating field '%s', the value got '%s' ,time parse throws an error '%s'", vType.Field(i).Name, valueStr, er.Error()), nil
					}
				}
			}
		} else {
			if IsInt(tagValue) {
				_, er := strconv.Atoi(valueStr)
				if er != nil {
					return false, vType.Field(i).Name + "format required int but got " + valueStr, nil
				}
			} else if IsFloat(tagValue) {
				_, er := strconv.ParseFloat(valueStr, 64)
				if er != nil {
					return false, vType.Field(i).Name + "format required float but got " + valueStr, nil
				}

			} else if in(tagValue , "time.Time", "time.time") {
				//"2006/1/2 15:04:05"
				_, er := time.ParseInLocation("2006/1/2 15:04:05", valueStr, time.Local)
				if er != nil {
					return false, fmt.Sprintf("while validating field '%s', the value got '%s' ,time parse throws an error '%s'", vType.Field(i).Name, valueStr, er.Error()), nil
				}
			}
		}
		continue
	}
	return true, "success", nil
}

//  validate the tag whose key named 'validate'
// the same as FormatCheck,but sounds more specific
func (checker *Checker) Validate(input interface{}) (bool, string, error) {
	return checker.FormatCheck(input)
}

// validate a struct methods whose method name ends with '"SVValidate"+<typ flags>'.
// typ is short for 'type' to avoid the built-in word type.
// what the string 'typ' stands for?
// assume you are going to validate a struct User{Username string , Age int},
// when creating a new user, 'username' field is required, however when updating a user,'username' is optional.
// typ helps to tell which method validates on which case.
// For example:
// func (o O) UserSVValidateSVBCreate()(bool,string,error){
//}
// SV means 'Super Validate' ,which tells the checker this method will be checked when call ValidateMethods.
// SVB means 'Super Validate Begin', which tells the checker to identify the typ flags.
// in the example above, 'typ' is 'Create',and this method only works for 'checker.ValidateMethods(o, 'Create')'
// func (o O) UserSVValidateSVBCreateSVSUpdate()(bool,string,error){
//}
// SV means 'Super Validate' ,which tells the checker this method will be checked when call ValidateMethods.
// SVB means 'Super Validate Begin', which tells the checker to identify the typ flags.
// SVS meas 'Super Validate Separate', which takes typ flags apart.
// in the example above, this method only works for 'checker.ValidateMethods(o, 'Create', 'Update')'.
// if typ is not set, all methods which end with 'Validate' or 'SVValidate' will be well checked
func (checker *Checker) ValidateMethods(input interface{}, typ ...string) (bool, string, error) {
	vType := reflect.TypeOf(input)
	vValue := reflect.ValueOf(input)
	var info string
	var methodName string

	var results []reflect.Value

	for i := 0; i < vType.NumMethod(); i++ {
		methodName = vType.Method(i).Name

		// UserValidate,UserSVValidate
		if strings.HasSuffix(methodName, "Validate") || strings.HasSuffix(methodName, "SVValidate") {
			// all cases will validate methods end with 'Validate' or 'SVValidate'
			results = vValue.Method(i).Call(nil)
			if len(results) != 3 {
				info = fmt.Sprintf("while validating method[%d],named '%s',illegal return values,want 3(bool,string,error) but got %d(%s)", i, methodName, len(results), valueListByType(results))
				return false, info, errors.New(info)
			}
			var er error
			ok, msg := results[0].Bool(), results[1].String()
			if results[2].IsNil() {
				er = nil
			} else {
				er = results[2].Interface().(error)
			}
			if ok {
				continue
			} else {
				return ok, msg, er
			}
		} else {
			if !(strings.Contains(methodName, "Validate") || strings.Contains(methodName, "SVValidate")) {
				// UserVx
				continue
			} else {
				// UserValidateXXX, UserSVValidateXXX
				if !(strings.Contains(methodName, "ValidateSVB") || strings.Contains(methodName, "SVValidateSVB")) {
					// UserValidateDrm
					continue
				} else {
					// UserValidateSVBXXX
					// validate by typ flags
					var declared []string
					j := strings.Index(methodName, "SVB")
					declared = strings.Split(methodName[j+len("SVB"):], "SVS")
					if hits(declared, typ) {
						results = vValue.Method(i).Call(nil)
						if len(results) != 3 {
							info = fmt.Sprintf("while validating method[%d],named '%s',illegal return values,want 3(bool,string,error) but got %d(%s)", i, methodName, len(results), valueListByType(results))
							return false, info, errors.New(info)
						}
						var er error
						ok, msg := results[0].Bool(), results[1].String()
						if results[2].IsNil() {
							er = nil
						} else {
							er = results[2].Interface().(error)
						}
						if ok {
							continue
						} else {
							return ok, msg, er
						}
					} else {
						continue
					}
				}
			}

		}
	}
	return true, "success", nil
}

// when input a []reflect.Value{false, 5, 'example'}
// returns 'bool,int,string'
func valueListByType(r []reflect.Value) string {
	typs := make([]string, 0)
	for _, v := range r {
		typs = append(typs, reflect.TypeOf(v.Interface()).String())
	}
	return strings.Join(typs, ",")
}

// whether two array can hits each other.
// ["1","2","3"] hits ["2","3"]
// ["1", "2", "3"] doesn't hit ["4"]
func hits(arr1 []string, arr2 []string) bool {
	for _, v1 := range arr1 {
		for _, v2 := range arr2 {
			if strings.ToLower(v1) == strings.ToLower(v2) {
				return true
			}
		}
	}
	return false
}

func checkRegex(input string, regex *regexp.Regexp) bool {
	return regex.MatchString(input)
}

func whetherLowerCase(tagValue string) bool {
	if strings.HasPrefix(tagValue, "regex") && strings.Contains(tagValue, "^") && strings.Contains(tagValue, "$") {
		return false
	}
	return true
}

func rollingCheck(checker *Checker, valueStr string, tagValue string, symbol string) (bool, error) {
	subStrings := strings.Split(tagValue, symbol)
	for i, v := range subStrings {
		if !checker.IsContainKey(v) {
			return false, errors.New(fmt.Sprintf("regex rule '%s' undefined", v))
		}
		if checker.IsBuilderContainKey(v) {
			//fmt.Println(fmt.Sprintf("'%s' id defined in added pool")

			if !checkRegex(valueStr, checker.ruler.RegexBuilder[v]) {
				//fmt.Println(fmt.Sprintf("'%s' match fail", v))
				return false, nil
			} else {
				if symbol == "|" {
					return true, nil
				}
				continue
			}
		}
		if !checkRegex(valueStr, checker.GetDefaultBuilder()[v]) {
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

// check an input string value by a raw regex string
func (checker *Checker) Check(input string, regex string) (bool, error) {
	r, er := regexp.Compile(regex)
	if er != nil {
		return false, er
	}
	return r.MatchString(input), nil
}

// check an input string value by the compiled regex object from the checker's default and added pool
func (checker *Checker) CheckFromPool(input string, key string) (bool, error) {
	key = strings.ToLower(key)
	if !checker.IsContainKey(key) {
		return false, errors.New(fmt.Sprintf("key '%s' not found in any of default or added regex pool", key))
	}

	r, ok := checker.ruler.RegexBuilder[key]
	if !ok {
		return checker.ruler.defaultRegexBuilder[key].MatchString(input), nil
	}
	return r.MatchString(input), nil
}

// This method is used for the case when dest struct has no 'validate' tag and tag value.
// Tag name can be specific,like 'json'.
//
// tagKey value should be single `json:""username` or 'username' as its first word splitting by ','
// these tags are ok `json:"username"`, `json:"username,x"`, `json:"username,x,y,z,n,m,q"
//
// type User struct{
//     Username string `json:"username"`
// }
//
// user := User{ Username:"LiLei" }
// checker := SuperChecker.GetChecker()
// ok,msg,e:=checker.ValidateByTagKeyAndMapValue(user, "json", map[string]string{ "username": "regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$")}
// fmt.Println(ok, msg, e)
func (checker *Checker) ValidateByTagKeyAndMapValue(dest interface{}, tagKey string, tags map[string]string) (bool, string, error) {
	dest = reflect.Indirect(reflect.ValueOf(dest)).Interface()

	vtype := reflect.TypeOf(dest)
	vvalue := reflect.ValueOf(dest)

	// for User struct,
	// "UserName" is typeName,
	// "json" is tagKey,
	// "username" is tagKeyValue,
	// "regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$" is tagValue
	var typeName string
	var tagKeyValue string
	var tagValue string

	var ok bool
	var msg string
	var e error

	for i := 0; i < vvalue.NumField(); i ++ {
		tagKeyValue = strings.Split(vtype.Field(i).Tag.Get(tagKey),",")[0]
		if tagKeyValue == "" || tagKeyValue == "-" {
			continue
		}
		tagValue, ok =tags[tagKeyValue]
		if !ok {
			continue
		}

		typeName = vtype.Field(i).Name
        ok,msg, e =checker.ValidateOne(typeName, vvalue.Field(i).Interface(), tagValue)
        if ok && e==nil{
        	continue
		}
        return ok,msg,e
	}
    return true, "success", nil
}

func (checker *Checker) ValidateOne(typeName string, value interface{}, tagValue string) (bool, string, error) {
	if whetherLowerCase(tagValue) {
		tagValue = strings.ToLower(tagValue)
	}

	// `validate:""` `validate:"-"` will be ignored
	if tagValue == "" || tagValue == "-" {
		return true, "success", nil
	}

	// the empty value will be ignore if no 'notnull' flag
	if !strings.Contains(strings.ToLower(tagValue), "notnull") {
		if ToString(value) == "" {
			return true, "success", nil
		}
		if whetherLowerCase(tagValue) {
			tagValue = strings.ToLower(tagValue)
		}

		// `validate:""` `validate:"-"` will be ignored
		if tagValue == "" || tagValue == "-" {
			return true, "success", nil
		}

		// the empty value will be ignore if no 'notnull' flag
		if !strings.Contains(strings.ToLower(tagValue), "notnull") {
			var zeroTime time.Time
			valueStr := ToString(value)
			if valueStr == "" || valueStr == zeroTime.Format("2006-01-02 15:04:05") {
				return true, "success", nil
			}
		}

		tmp := strings.Split(tagValue, ",")
		flag := strings.ToLower(tmp[0])
		if !strings.Contains(strings.Join(flagRange, " "), flag) {
			return false,
				fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", typeName, flag, flagRange),
				errors.New(fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", typeName, flag, flagRange),
				)
		}
	}

	tmp := strings.Split(tagValue, ",")
	flag := strings.ToLower(tmp[0])
	if !strings.Contains(strings.Join(flagRange, " "), flag) {
		return false,
			fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", typeName, flag, flagRange),
			errors.New(fmt.Sprintf("while validating field '%s',flag '%s' is not contained in the flagRange '%v'", typeName, flag, flagRange),
			)
	}

	var rule string
	if len(tmp) > 1 {
		rule = strings.Join(tmp[1:], ",")
	}

	var valueStr string
	// flag range/in validate
	if flag == "range" || flag == "in" {
		if !strings.HasPrefix(rule, "[") || !strings.HasSuffix(rule, "]") {
			return false,
				fmt.Sprintf("field '%s' range/in flag must have its rule format like '[x,x,x,x,x] but got '%s'", typeName, rule),
				errors.New(fmt.Sprintf("field '%s' range/in flag must have its rule format like '[x,x,x,x,x] but got '%s'", typeName, rule))
		}
		valueStr = ToString(value)
		arr := strings.Split(rule[1:len(rule)-1], ",")
		for _, v := range arr {
			if v == valueStr {
				return true, "success", nil
			}
		}
		return false, fmt.Sprintf("field '%s' required in '%s' but got '%s'", typeName, rule, valueStr), nil
	}

	// regex validate
	// regex validate is used to replace superChecker tag
	// Name string `superChecker:"key"`  <==> Name string `validate:"regex,key"`
	if flag == "regex" {
		if whetherLowerCase(tagValue) {
			rule = strings.ToLower(rule)
		}
		valueStr = ToString(value)
		// rule is raw regex
		if strings.HasPrefix(rule, `^`) && strings.HasSuffix(rule, `$`) {
			ok, er := checker.Check(valueStr, rule)
			if er != nil {
				return false,
					fmt.Sprintf("while validating field '%s' regex '%s' throws an error '%s'", typeName, strconv.QuoteToASCII(rule), er.Error()),
					errors.New(fmt.Sprintf("while validating field '%s' regex '%s' throws an error '%s'", typeName, strconv.QuoteToASCII(rule), er.Error()))
			}
			if !ok {
				return false,
					fmt.Sprintf("while validating field '%s' regex '%s' but got unmatched value '%s'", typeName, strconv.QuoteToASCII(rule), valueStr),
					nil
			}
			return true, "success", nil
		} else {
			if strings.Contains(rule, "|") {
				// rule formated like 'key1|key2|key3' which can be separated by '|'
				rules := strings.Split(rule, "|")
				for j, v := range rules {
					ok, er := checker.CheckFromPool(valueStr, v)
					if er != nil {
						return false,
							fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", typeName, j, v, er.Error()),
							errors.New(fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", typeName, j, v, er.Error()))
					}
					if ok {
						return true, "success", nil
					}
					if j == len(rules)-1 {
						return false, fmt.Sprintf("while validating field '%s', regex key group %v all fail, got unmatched value '%s'", typeName, rules, valueStr), nil
					}
				}
			} else if strings.Contains(rule, ",") {
				// rule formated like 'key1,key2,key3' which can be separated by ','
				rules := strings.Split(rule, ",")
				for i, v := range rules {
					ok, er := checker.CheckFromPool(valueStr, v)
					if er != nil {
						return false,
							fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", typeName, i, v, er.Error()),
							errors.New(fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' throws an error '%s'", typeName, i, v, er.Error()))
					}
					if !ok {
						return false,
							fmt.Sprintf("while validating field '%s', regex group['%d'] regex pool key '%s' but got unmatched value '%s'", typeName, i, v, valueStr),
							nil
					}
				}
			} else {
				// rule is regarded as the key itself
				if !checker.IsContainKey(rule) {
					return false,
						fmt.Sprintf("while validating field '%s' regex key '%s' not found in any of default regex pool or add regex pool,you may use '%s' before using it", typeName, rule, "checker.AddRegex('key', '^raw regex$')"),
						errors.New(fmt.Sprintf("while validating field '%s' regex key '%s' not found in any of default regex pool or add regex pool,you may use '%s' before using it", typeName, rule, "checker.AddRegex('key', '^raw regex$')"))
				}
				ok, er := checker.CheckFromPool(valueStr, rule)
				if er != nil {
					return false,
						fmt.Sprintf("while validating field '%s' regex pool key '%s' throws an error '%s'", typeName, rule, er.Error()),
						er
				}
				if !ok {
					return false,
						fmt.Sprintf("while validating field '%s' regex pool key '%s' but got unmatched value '%s'", typeName, rule, valueStr),
						nil
				}
			}
		}
		return true, "success", nil
	}

	var ok1, ok2, ok3, ok bool
	// type validate including:
	// func validate,
	// int,float,time validate
	_, ok1 = value.(time.Time)
	_, ok3 = value.(jsoncrack.Time)
	ok2 = strings.Split(tagValue, ",")[0] == "time.time"
	ok = ok1 || ok2 || ok3
	if ok && strings.Contains(tagValue, ",") && strings.Split(tagValue, ",")[1] != "" {
		valueStr = ToString(value, strings.Split(tagValue, ",")[1])
	} else {
		valueStr = ToString(value)
	}

	if valueStr == "undefined" || valueStr == "undefine" || valueStr == "" {
		return true, "success", nil
	}

	if strings.Contains(tagValue, ",") {
		tmp = strings.Split(tagValue, ",")
		tagValue = tmp[0]
		rule = strings.Join(tmp[1:], ",")
		if isFunc(tagValue) {
			if len(tmp) < 2 {
				return false,
					fmt.Sprintf("'%s' is validated as 'func', the tag 'validate' must has its tag value length more than 2,but got '%s' length is %d", typeName, tagValue, len(tmp)),
					errors.New(fmt.Sprintf("'%s' is validated as 'func', the tag 'validate' must has its tag value length more than 2,but got '%s' length %d", typeName, tagValue, len(tmp)))
			}

			if len(tmp) == 2 {
				if arr := strings.Split(rule, "|"); len(arr) > 1 {
					// validate:func,key1|key2|key3
					for j, r := range arr {
						if !checker.ContainFunc(r) {
							return false, fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", typeName, j, r),
								errors.New(fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", typeName, j, r))
						}
						ok, msg, er := checker.GetFunc(r).Value(value, typeName)
						if ok {
							return true, "success", nil
						}
						if j >= len(arr)-1 {
							return ok, msg, er
						}
						continue
					}
				} else {
					// validate:func,key1
					if !checker.ContainFunc(rule) {
						return false, fmt.Sprintf("'%s' func has not be added into func pool,use checker.AddFunc() to register", rule),
							errors.New(fmt.Sprintf("'%s' func has not be added into func pool,use checker.AddFunc() to register", rule))
					}
					ok, msg, er := checker.GetFunc(rule).Value(value, typeName)
					if ok {
						return true, "success", nil
					}
					return ok, msg, er
				}
			} else {
				// validate:func,key1,key2,key3,key4
				rules := strings.Split(rule, ",")
				for j, r := range rules {
					if !checker.ContainFunc(r) {
						return false, fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", typeName, j, r),
							errors.New(fmt.Sprintf("while validating field '%s', func group[%d] '%s' func has not be added into func pool,use checker.AddFunc() to register", typeName, j, r))
					}
					ok, msg, er := checker.GetFunc(r).Value(value, typeName)
					if !ok {
						return ok, msg, er
					}
					if j == len(rules)-1 {
						return true, "success", nil
					}
				}
			}

		} else if IsInt(tagValue) {
			if rule != "" {
				tmp2 := strings.Split(rule, ":")
				if len(tmp2) != 2 {
					return false, "", errors.New("notation requires number1:number2,but got " + rule)
				}
				v, er := strconv.Atoi(valueStr)
				if er != nil {
					return false, typeName + " format required int but got " + valueStr, nil
				}
				if tmp2[0] != "" {
					min, er := strconv.Atoi(tmp2[0])
					if er != nil {
						return false, "", errors.New(typeName + " notation rule required number:number but get " + tmp2[0])
					}
					if v < min {
						return false, typeName + " int value required bigger than " + tmp2[0] + " but get " + valueStr, nil
					}
				}
				if tmp2[1] != "" {
					max, er := strconv.Atoi(tmp2[1])
					if er != nil {
						return false, "", errors.New(typeName + " notation rule required number:number but get " + tmp2[1])
					}
					if v > max {
						return false, typeName + " int value required smaller than " + tmp2[1] + " but get " + valueStr, nil
					}
				}

			} else {
				_, er := strconv.Atoi(valueStr)
				if er != nil {
					return false, typeName + " format required int but got " + valueStr, nil
				}
			}
		} else if IsFloat(tagValue) {
			if rule != "" {
				tmp2 := strings.Split(rule, ":")
				if len(tmp2) != 2 {
					return false, "", errors.New(" notation requires float_number1:float_number2,but got " + rule)
				}
				v, er := strconv.ParseFloat(valueStr, 64)
				if er != nil {
					return false, typeName + " format required float but got " + valueStr, nil
				}
				if tmp2[0] != "" {
					min, er := strconv.ParseFloat(tmp2[0], 64)
					if er != nil {
						return false, "", errors.New(typeName + " notation rule required float_number:float_number but got " + tmp2[0])
					}
					if v < min {
						return false, typeName + " float value required bigger than" + tmp2[0] + " but got " + valueStr, nil
					}
				}
				if tmp2[1] != "" {
					max, er := strconv.ParseFloat(tmp2[1], 64)
					if er != nil {
						return false, "", errors.New(typeName + " notation rule required number:number but got " + tmp2[1])
					}
					if v > max {
						return false, typeName + " float value required smaller than " + tmp2[1] + " but got " + valueStr, nil
					}
				}
			} else {
				_, er := strconv.ParseFloat(valueStr, 64)
				if er != nil {
					return false, typeName + "format required float but got " + valueStr, nil
				}
			}
		} else if tagValue == "time.time" || tagValue=="time.Time" {
			//"2006/1/2 15:04:05"
			if rule != "" {
				_, er := time.ParseInLocation(rule, valueStr, time.Local)
				if er != nil {
					return false, fmt.Sprintf("while validating field '%s', time format requires %s but go %s", typeName, rule, valueStr), nil
				}
			} else {
				_, er := time.ParseInLocation("2006/1/2 15:04:05", valueStr, time.Local)
				if er != nil {
					return false, fmt.Sprintf("while validating field '%s', time format requires %s but go %s", typeName, rule, valueStr), nil
				}
			}
		}
	} else {
		if IsInt(tagValue) {
			_, er := strconv.Atoi(valueStr)
			if er != nil {
				return false, typeName + "format required int but got " + valueStr, nil
			}
		} else if IsFloat(tagValue) {
			_, er := strconv.ParseFloat(valueStr, 64)
			if er != nil {
				return false, typeName + "format required float but got " + valueStr, nil
			}

		} else if tagValue == "time.Time" {
			//"2006/1/2 15:04:05"
			_, er := time.ParseInLocation("2006/1/2 15:04:05", valueStr, time.Local)
			if er != nil {
				return false, fmt.Sprintf("while validating field '%s', the value got '%s' ,time parse throws an error '%s'", typeName, valueStr, er.Error()), nil
			}
		}
	}
	return true, "success", nil
}

// int type assertion
func IsInt(in string) bool {
	in = strings.ToLower(in)
	for _, v := range intTypes {
		if v == in {
			return true
		}
	}
	return false
}

// float type assertion
func IsFloat(in string) bool {
	in = strings.ToLower(in)
	for _, v := range floatTypes {
		if v == in {
			return true
		}
	}
	return false
}

// func type assertion
func isFunc(in string) bool {
	in = strings.ToLower(in)
	if in == "func" || in == "function" {
		return true
	}
	return false
}
