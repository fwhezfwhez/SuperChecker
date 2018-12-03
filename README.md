a validator and checker tool. validator works for validating whether the input data is valid, and superchecker works for checking its value by regex

# superchecker
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fwhezfwhez/SuperChecker)

## Example
```go
package main

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"strings"
	"superChecker"
	"time"
)

type User struct {
    // validate via func pool
	Introduce    string       `validate:"func,introduction"`

    // validate via validation tag
	InTime    time.Time       `validate:"time.Time,2006.01.2 15:04:05"`
	InTimeStr string          `validate:"time.Time,2006.01.2 15:04:05"`
	Age       int             `validate:"int,:140"`
	Salary    float64         `validate:"float"`
	Deci      decimal.Decimal `validate:"float,10:100"`

    // jump
	K        []int

    // check via regex
	Phone    string `superChecker:"mobilephone|telephone"`
	UserName string `superChecker:"userName" json:"userName" `
	Password string `superChecker:"password"`
	NotNull  string `superChecker:"notnull"`
	Text   string //`superChecker:"length,chineseOnly,notNull"`
	Number int64 `superChecker:"mobilephone"`
}

func main() {
    // build a testing data
	user := User{
		Introduce: "I am flying",
		InTime:    time.Now(),
		InTimeStr: time.Now().Format("2006.01-2 15:04:05"),
		UserName:  "d",
		Password:  "a1dfdasfsdf",
		Phone:     "13875847584",
		Text:      "undefined",
		Age:       130,
		Salary:    3000.9,
		Number:    18970937633,
		NotNull:   "3",
		Deci:      decimal.NewFromFloat(11.9),
	}
    // init a checker object
	checker := superChecker.GetChecker()

    // add regex rule into regex pool
	checker.AddRegex("passWoRd", "^[\\s\\S]{6,}$")
	checker.AddRegex("length", "^[\\s\\S]{0,20}$")
	checker.AddRegex("chineseOnly", "^[\u4E00-\u9FA5]*$")

    // add func into func pool
	checker.AddFunc(func(in interface{})(bool,string,error){
		v,ok := in.(string)
		if !ok {
			return false, "assertion error,in is not a string type", errors.New("assertion error,in is not a string type")
		}
		// deal with v
		// length limit
		if len(v) >1000 {
			return false, fmt.Sprintf("max len is 1000,but got %d", len(v)), nil
		}
		// abuse words limit
		if strings.Contains(v,"fuck") {
			return false, fmt.Sprintf("'%s' contains bad words '%s'", v, "fuck"), nil
		}
		return true,"success",nil
	}, "introduction")

    // add regex rule into default pool
	checker.AddDefaultRegex("chineseOnly", "^[\u4E00-\u9FA5]*$")
	fmt.Println("-------Check(in string,rule string)---------")

    // regex checks a single string data
	ok, er := checker.Check("10000124", "^[0-9]{8}$")
	fmt.Println(ok, er)

	// check those fields whose tag is 'superChecker'
	result, msg, err := checker.SuperCheck(user)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("-------SuperCheck(in interface{})-------")

	fmt.Println("匹配结果:", result, "信息:", msg)

	fmt.Println("-------FormatCheck(in interface{})/Validate(in interface{})-------")
	// validate those fields whose tag is 'validate'
	ok, msg, er = checker.FormatCheck(user)
	//ok, msg, er = checker.Validate(user)

	if er != nil {
		log.Println(er.Error())
	}
	fmt.Println("格式验证结果:", ok, "msg:", msg)

	fmt.Println("-------listAll----------")
	checker.ListAll()
}

```

## How to specific superchecker tag?
**superChecker**:
The tag value is the key added by **AddRegex** or **AddDefaultRegex**, while the former one is the added pool which has higher privalage than the latter one(when both of them has a key 'password', than use the regex in added pool).

```go
   type User struct {
	Password string `superChecker:"password"
   }
   ...
   checker.AddRegex("password", "^[\\s\\S]{6,}$")
```

When a field will fit several regex rules, use it like
```go
   type User struct {
	Phone string `superChecker:"phone|mobilePhone"
	Introduction string `superChecker:"length,noAbuse,noChinese"`
   }
   checker.AddRegex("phone",  "^[0-9]{8}$")
   checker.AddRegex("mobilePhone","^1[0-9]{10}$")
   ...
```
`key1|key2|key3` means the field(Phone) should fit one of keys(phone,mobilePhone), the or logic.
`key1,key2,key3` means the field(Introduction) should both fit all of the keys(length,noAbuse,noChinese)

**I'm sorry but checker doesn't support '|' and ',' mixed like `key1,key2|key3`, also doesn't support quoted like 'key1,key2,(key3,key4)'. Soon the checker will give its solutions to this situation**

## How to specific validate tag?
**validate**:
The tag value consists of two parts, type and rule(key).
type and rule used like:
```go
type User struct{
    Age int `validate:"int,0:200"`  // age should be integer and between 0 and 200
	// Age int `validate:"int,:200"`  // age should be integer and less than 200
	// Age int `validate:"int,0:"`  // age should be integer and bigger than 0

    Salary float64 `validate:"float,0:1000000000"`  // Salary  should be float type(float32,float64) and between 0 and 1000000000
	// Salary float64 `validate:"float,:1000000000"`  // Salary  should be float type(float32,float64) and less than 1000000000
	// Salary float64 `validate:"float,0:"`  // Salary  should be float type(float32,float64) and bigger than 0

	// InTime    time.Time       `validate:"time.Time"`// golang support deliver the origin time type ,it's good to use time.Time field to bind data
	// if insist on using string type to bind time data,use it like:
    InTimeStr string          `validate:"time.Time,2006.01.2 15:04:05"` // InTimeStr should fit the format '2006.01.2 15:04:05'
}
```
The tag value like 'time.Time','int','float' is the type ,and the latter string words is its rule,like '0:200'.
**int means int type ,it's ok to write like:**
```go
...
Age `validate:"int,:120"`
//Age `validate:"int8,:120"`
//Age `validate:"int16,:120"`
//Age `validate:"int32,:120"`
//Age `validate:"int64,:120"`
```
**so does float types**
```go
...
Salary `validate:"float,:120"`
//Salary `validate:"float32,:120"`
//Salary `validate:"float64,:120"`
//Salary `validate:"decimal,:120"`
```
## How to design a function to validate data?
```go
...
type User struct {
	Introduce    string       `validate:"func,introduction"`
}
...
checker.AddFunc(func(in interface{})(bool,string,error){
		v,ok := in.(string)
		if !ok {
			return false, "assertion error,in is not a string type", errors.New("assertion error,in is not a string type")
		}
		// deal with v
		// length limit
		if len(v) >1000 {
			return false, fmt.Sprintf("max len is 1000,but got %d", len(v)), nil
		}
		// abuse words limit
		if strings.Contains(v,"fuck") {
			return false, fmt.Sprintf("'%s' contains bad words '%s'", v, "fuck"), nil
		}
		return true,"success",nil
	}, "introduction")
...
```
**More tips on developing and if you want to help contribute,please fork and pull request. More issues please hand in in issue part.3q**
