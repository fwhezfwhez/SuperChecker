a validator and checker tool. validator works for validating whether the input data is valid

# superchecker
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fwhezfwhez/SuperChecker)
[![Build Status]( https://www.travis-ci.org/fwhezfwhez/SuperChecker.svg?branch=master)]( https://www.travis-ci.org/fwhezfwhez/SuperChecker)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [superchecker](#superchecker)
  - [1. Start](#1-start)
  - [2. Tips](#2-tips)
      - [2.1 The tag 'superChecker' has stopped developing](#21-the-tag-superchecker-has-stopped-developing)
      - [2.2 There are some built-in regex key, case-not-sensitive](#22-there-are-some-built-in-regex-key-case-not-sensitive)
      - [2.3 regex key combining](#23-regex-key-combining)
  - [3. Example](#3-example)
      - [3.1 superChecker tag example(this tag is stoped developing,and replaced by validate, the old usage is still access)](#31-superchecker-tag-examplethis-tag-is-stoped-developingand-replaced-by-validate-the-old-usage-is-still-access)
      - [3.2 validate tag example](#32-validate-tag-example)
  - [4. FAQ](#4-faq)
      - [4.1 How to specific superchecker tag?](#41-how-to-specific-superchecker-tag)
      - [4.2 How to specific validate tag?](#42-how-to-specific-validate-tag)
      - [4.3 How to design a function to validate data?](#43-how-to-design-a-function-to-validate-data)
      - [4.4 How to validate a model by its method and the model as receiver](#44-how-to-validate-a-model-by-its-method-and-the-model-as-receiver)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



## 1. Start
`go get github.com/fwhezfwhez/SuperChecker`

## 2. Tips
#### 2.1 The tag 'superChecker' has stopped developing
the tag 'superChecker' **is abandoned** and replaced by tag 'validate'. Old usages are remained, function is still access if you insist on using this.More usages refer to Example
```go
type U struct{
Username string `superChecker:"username"`
}
```

#### 2.2 There are some built-in regex key, case-not-sensitive
| key | regex | desc |
|:--------|:--------|:--- |
| chineseOnly | ^[\u4E00-\u9FA5]*$ | only chinese format |
| notNull | ^[\\s\\S]+$ | not empty |
| username | ^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$ | chinese,a-z,A-Z,0-9,_,. len 0-40|
| number | ^[0-9]+$ | number more than one bit |
| decimal | ^\\d+\\.[0-9]+$ | number |
| mobilePhone | ^1[0-9]{10}$ | 10 bit mobile phone |
| telephone | ^[0-9]{8}$ | 8 bit telephone |

#### 2.3 regex key combining
**logic 'and' and 'or'**
`key1,key2,key3 ...` the value should fit all of the keys
`key1|key2|key3 ...` the value should fit one of the keys

** not supported:**
- `key1,key2,key3|key4,key5|key6` use ',' '|' together not supported
- `key1,key2,(key3,key4)` use '()' is not supported
```go
type U struct{
Username string `superChecker:"username,notNull"`
Username2 string `validate:"regex,username,notNull"`
Username3 string `superChecker:"username|notNull"`
Username4 string `validate:"regex,username|notNull"`
}
```

## 3. Example
#### 3.1 superChecker tag example(this tag is stoped developing,and replaced by validate, the old usage is still access)
```go
// Tips:
// 'superChecker' tag has been stopped developing, the old functions will be remained.
// 'validate' tag is to replace the old ussages:
// "Name string `superChecker:"key1,key2"`" equals to "Name string `validate:"regex,key1,key2"`"

package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/fwhezfwhez/superChecker"
)

type Animal struct{
	Name string `superChecker:"username"`
	Count int `superChecker:"positive"`
	Price decimal.Decimal `superChecker:"positive"`
}

func main() {
	animal := Animal{
		Name:"beibei",
		Count: 1000,
		Price: decimal.NewFromFloat(100000),
	}

	checker := superChecker.GetChecker()
	checker.AddRegex("username","^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$")
	checker.AddRegex("positive","^[0-9.]+$")
	ok,msg,er := checker.SuperCheck(animal)
	if er!=nil {
		fmt.Println(fmt.Sprintf("got an error: '%s'", er.Error()))
		return
	}
	if !ok {
		fmt.Println(fmt.Sprintf("fail because of : '%s'", msg))
		return
	}
	fmt.Println("success")
}
```

#### 3.2 validate tag example
```go
package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/fwhezfwhez/SuperChecker"
	"time"
)

type Order struct {
	// TIME
	CreatedAt time.Time `validate:"time.time"`
	UpdatedAt string    `validate:"time.time,2006/01/02 15:04:05"`

	// INT
	Count    int `validate:"int,0:200"`
	MaxCount int `validate:"int,:200"`
	MinCount int `validate:"int,10:"`
	Count2   int `validate:"int64,0:200"`

	// FLOAT
	RewardRate    float64         `validate:"float,0:0.4"`
	MaxRewardRate float64         `validate:"float,:0.4"`
	MinRewradRate float64         `validate:"float,0:"`
	RewardRate2   float64         `validate:"float64,0:0.4"`
	RewardRate3   decimal.Decimal `validate:"decimal,0:0.4"`

	// REGEX
	OrderUsername  string `validate:"regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$"`
	OrderUsername2 string `validate:"regex,username"`

	// RANGE,IN
	OrderStatus     int    `validate:"range,[1,2,3,4]"`
	OrderStatusName string `validate:"in,[unpaid,paid,closed]"`

	// FUNC, FUNCTION
	MailTypeCheckBox  string `validate:"func,inAndLength,lengthMoreThan3"`
	MailTypeCheckBox2 string `validate:"function,lengthLessThan3|inAndLength"`
}
func (o Order) XXSVValidateSVBCreate()(bool,string,error){
	return true,"xxsvcreate wrong",nil
}
func (o Order) XXValidate()(bool,string,error){
	return true,"xxv wrong",nil
}
func (o Order) XXSVValidate()(bool,string,error){
	return true,"xxsv wrong",nil
}

func (o Order) XXValidateSVBCreate()(bool,string,error){
	return true,"xxcreate wrong",nil
}



func (o Order) XXValidateSVBCreateSVSUpdate()(bool,string,error){
	return false,"xxsvcreateupdate wrong",nil
}
func (o Order) XXSVValidateSVBCreateSVSUpdate()(bool,string,error){
	return true,"xxsvcreateupdate wrong",nil
}

func main() {
	order := Order{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Format("2006/01/02 15:04:05"),

		Count:    200,
		MaxCount: 90,
		MinCount: 10,
		Count2:   100,

		RewardRate:    0.4,
		MaxRewardRate: 0.3,
		MinRewradRate: 0.1,
		RewardRate2:   0.1,
		RewardRate3:   decimal.NewFromFloat(0.4),

		OrderUsername:  "superCheckerValidate",
		OrderUsername2: "superCheckerValidate",

		OrderStatus:     3,
		OrderStatusName: "closed",

		MailTypeCheckBox:  "midMail",
		MailTypeCheckBox2: "midMail",
	}

	checker := superChecker.GetChecker()
	checker.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		v := superChecker.ToString(in)
		maxLength := 7
		if len(v) > maxLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' over length,want %d ,but got %d", fieldName, "inAndLength", maxLength, len(v)), nil
		}
		vrange := []string{"midMail", "shenMail", "yundaMail"}
		for _, value := range vrange {
			if value == v {
				return true, "success", nil
			}
		}
		return false, fmt.Sprintf("while validating field '%s', rule key '%s',  value '%s' not in '%v'", fieldName, "inAndLength", v, vrange), nil
	}, "inAndLength")
	checker.AddFunc(func(in interface{}, fieldName string)(bool, string, error){
		v := superChecker.ToString(in)
		minLength := 3
		if len(v) < minLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' too short length,want %d ,but got %d", fieldName, "inAndLength", minLength, len(v)), nil
		}
		return true, "success", nil
	},"lengthmorethan3")

	checker.AddFunc(func(in interface{}, fieldName string)(bool, string, error){
		v := superChecker.ToString(in)
		maxLength := 3
		if len(v) > maxLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' too short length,want %d ,but got %d", fieldName, "inAndLength", maxLength, len(v)), nil
		}
		return true, "success", nil
	},"lengthlessthan3")

	ok, msg, er := checker.Validate(order)
	if er != nil {
		fmt.Println(fmt.Sprintf("got an error, '%s'", er.Error()))
		return
	}
	if !ok {
		fmt.Println(fmt.Sprintf("validate fail because of '%s'", msg))
		return
	}


	// ioc, inverse of control
	// validate to combine as receiver to the dest struct
	ok, msg, er = checker.ValidateMethods(order,"create","update")
	if er != nil {
		fmt.Println(fmt.Sprintf("got an error, '%s'", er.Error()))
		return
	}
	if !ok {
		fmt.Println(fmt.Sprintf("validate fail because of '%s'", msg))
		return
	}
	fmt.Println("success")
}

```
## 4. FAQ
#### 4.1 How to specific superchecker tag?
**superChecker**(**abandoned**):
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

#### 4.2 How to specific validate tag?
**validate**:
The tag value consists of two parts, **type** and **rule(key)**.Let's see it in table below

| type | rule/key | desc |
|:--------| :--------| :---- |
| int/int8/int16/int32/int64, float/float32/float64/decimal | 0:100 | int value between 0 and 100, containing 0 and 100 |
| int/int8/int16/int32/int64 | 0: | int value bigger than 0, containing 0 |
| time.time | 2006/01/02 15:04:05| received time value should be formated like yyyy/MM/dd HH:mm:ss, struct type is string or time.Time or jsoncrack.Time, if is jsoncrack.Time or time.Time, the rule can be ignored like "CreatedAt time.Time \`validate:"time.time"\`" or just ignore tag "CreatedAt time.Time"|
| regex | ^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$ | value should be pass this regex validate, value should be string-able types(int,float,time,fmt.Stringer) |
| regex | username | value should pass an in-built or after-added regex key 'username', value should be string-able types(int,float,time,fmt.Stringer), make sure using "checker.AddRegex("username","^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$") before you use this key" |
| func/function | lengthlt10| value should be less than length by 10, make sure using "checker.AddFunc()" |
| range/in | [1,2,3,4]/[paid,unpaid,retreat] | value should be contained in the list |

type and rule used like:
```go
	type User struct {
		Age int `validate:"int,0:200"` // age should be integer and between 0 and 200
		// Age int `validate:"int,:200"`  // age should be integer and less than 200
		// Age int `validate:"int,0:"`  // age should be integer and bigger than 0

		Salary float64 `validate:"float,0:1000000000"` // Salary  should be float type(float32,float64) and between 0 and 1000000000
		// Salary float64 `validate:"float,:1000000000"`  // Salary  should be float type(float32,float64) and less than 1000000000
		// Salary float64 `validate:"float,0:"`  // Salary  should be float type(float32,float64) and bigger than 0

		// InTime    time.Time       `validate:"time.Time"`// golang support deliver the origin time type ,it's good to use time.Time field to bind data
		// if insist on using string type to bind time data,use it like:
		InTimeStr string `validate:"time.Time,2006.01.2 15:04:05"` // InTimeStr should fit the format '2006.01.2 15:04:05'
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
#### 4.3 How to design a function to validate data?
** normal **: `addFunc(f func(in interface{}, fieldName string)(bool, string, error), key)`
```go
...
type User struct {
	Introduce    string       `validate:"func,introduction"`
}
...
	checker.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		v, ok := in.(string)
		if !ok {
			return false, "assertion error,in is not a string type", errors.New("assertion error,in is not a string type")
		}
		// deal with v
		// length limit
		if len(v) > 1000 {
			return false, fmt.Sprintf("max len is 1000,but got %d", len(v)), nil
		}
		// abuse words limit
		if strings.Contains(v, "fuck") {
			return false, fmt.Sprintf("'%s' contains bad words '%s'", v, "fuck"), nil
		}
		return true, "success", nil
	}, "introduction")
...
```
**with path**: `addFunc(f func(in interface{}, fieldName string)(bool, string, error), key string, path string)`
the path is combined with a function where it's declared, used to help developer to locate where it is well set.

```go
...
type User struct {
	Introduce    string       `validate:"func,introduction"`
}
...
    _, file, line, _ := runtime.Caller(1)
	path = fmt.Sprintf("%s:%d", file, line)
	checker.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		v, ok := in.(string)
		if !ok {
			return false, "assertion error,in is not a string type", errors.New("assertion error,in is not a string type")
		}
		// deal with v
		// length limit
		if len(v) > 1000 {
			return false, fmt.Sprintf("max len is 1000,but got %d", len(v)), nil
		}
		// abuse words limit
		if strings.Contains(v, "fuck") {
			return false, fmt.Sprintf("'%s' contains bad words '%s'", v, "fuck"), nil
		}
		return true, "success", nil
	}, "introduction", path)
...
```

#### 4.4 How to validate a model by its method and the model as receiver
**Some concepts**:

| concept | short for what | example | declared typs(case-not-sensitive) | desc |
|:--------|:--------| :----| :----- |
| SVValidate | super valudate | `func (o Object) ObjectSVValidate()(bool,string,error)` | all | to declare the method is marked to be validated by `checker.ValidateMethods(o)` , sv can be ignored, `ObjectSVValidate` equals to `ObjectValidate`|
| SVB | super valudate begin | `func (o Object) ObjectSVValidateSVBCreate()(bool,string,error)` | create | to declare a marked method begin spot,after 'SVB' is the `typs` `checker.ValidateMethods(o, "create")` |
| SVS | super valudate seperate | `func (o Object) ObjectSVValidateSVBCreateSVSUpdate()(bool,string,error)` | create, update| to seperate 'typs' after 'SVB',`checker.ValidateMethods(o, "create", "update")` |

**details:**
- **typs choose which methods to be validate, they follow the HIT rule,** `[create, update]` **hits** `[ValidateSVBCreate, SVValidateSVBCreateSVSUpdate, Validate, ValidateSVBUpdate]`
- **typs are case-not-sensitive**

**example**:
```go
type O struct {
	Username string
}

// v1
func (o O) OLengthValidate() (bool, string, error) {
	if o.Username > 5 && o.Username < 100 {
		return true, "success", nil
	}
	return false, "length should be between 5 and 100", nil
}

// v2
func (o O) OValidateSVBCreate() (bool, string, error) {
	if o.Username != "" {
		return true, "success", nil
	}
	return false, "length should be between 5 and 100", nil
}

// v3
func (o O) OValidateSVBUpdate() (bool, string, error) {
	if o.Username == "admin" {
		return false, "admin should not be updated", nil
	}
	return true, "success", nil
}

// v4
func (o O) OValidateSVBUpdateSVSCreate() (bool, string, error) {
	if o.Username == "admin" {
		return false, "admin should not be updated", nil
	}
	return true, "success", nil
}
func main() {
	...
	// o:=O{Username:"he"}
	o := O{Username: "hellworld"}
	// v1 will be validated
	ok, msg, e := checker.ValidateMethods(o)

	// v1,v2,v4 will be validated
	ok, msg, e = checker.ValidateMethods(o, "create")

	// v1,v3,v4 will be validated
	ok, msg, e := checker.ValidateMethods(o, "update")

	// v1,v2,v3,v4 will all be validated
	ok, msg, e := checker.ValidateMethods(o, "update", "create")

	if e != nil {
		handle(e)
	}
	fmt.Println(ok, msg)
}
```

**More tips on developing and if you want to help contribute,please fork and pull request. More issues please hand in in issue part.3q**
