a validator and checker tool. validator works for validating whether the input data is valid, and superchecker works for checking its value by regex

# superchecker
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fwhezfwhez/SuperChecker)

## Example
### superChecker tag example
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

### validate tag example
```go
package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"superChecker"
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

	// RANGE, IN
	OrderStatus     int    `validate:"range,[1,2,3,4]"`
	OrderStatusName string `validate:"in,[unpaid,paid,closed]"`

	// FUNCTION, FUNC
	MailTypeCheckBox  string `validate:"func,inAndLength"`
	MailTypeCheckBox2 string `validate:"function,inAndLength"`
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
		MailTypeCheckBox2: "",
	}

	checker := superChecker.GetChecker()
	checker.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		v := superChecker.ToString(in)
		maxLength := 3
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
	ok, msg, er := checker.Validate(order)
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
checker.AddFunc(func(in interface{}, fieldName string)(bool,string,error){
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
