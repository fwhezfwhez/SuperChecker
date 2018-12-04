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
	// regex key validate
	Phone2 string `validate:"regex,mobilephone|telephone"`
	// regex raw validate
	Phone3 string `validate:"regex,^[0-9]{8}$"`

	// func validate
	Introduce    string       `validate:"func,introduction"`

	// type time.time validate
	InTime    time.Time       `validate:"time.Time"`
	InTimeStr string          `validate:"time.Time,2006.01.2 15:04:05"`

	// type int validate
	Age       int             `validate:"int,:140"`
	// type float validate
	Salary    float64         `validate:"float"`
	Deci      decimal.Decimal `validate:"float,10:100"`

	K        []int
	// superChecker
	//`superChecker:"mobilephone|telephone"` equals to `validate:"regex,mobilephone|telephone"`
	Phone    string `superChecker:"mobilephone|telephone"`


	UserName string `superChecker:"userName" json:"userName" `
	Password string `superChecker:"password"`
	NotNull  string `superChecker:"notnull"`

	Text   string //`superChecker:"length,chineseOnly,notNull"`
	Number int64 `superChecker:"mobilephone"`
}

func main() {
	user := User{
		Introduce: "I am flying",
		InTime:    time.Now(),
		InTimeStr: time.Now().Format("2006.01-2 15:04:05"),
		UserName:  "d",
		Password:  "a1dfdasfsdf",
		Phone:     "13875847584",
		Phone2: "",

		Text:      "undefined",
		Age:       130,
		Salary:    3000.9,
		Number:    18970937633,
		NotNull:   "3",
		Deci:      decimal.NewFromFloat(11.9),
	}
	checker := superChecker.GetChecker()
	checker.AddRegex("passWoRd", "^[\\s\\S]{6,}$")
	checker.AddRegex("length", "^[\\s\\S]{0,20}$")
	checker.AddRegex("chineseOnly", "^[\u4E00-\u9FA5]*$")

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

	checker.AddDefaultRegex("chineseOnly", "^[\u4E00-\u9FA5]*$")
	fmt.Println("-------Check(in string,rule string)---------------------------------------------------------------------------------------------")
	ok, er := checker.Check("10000124", "^[0-9]{8}$")
	fmt.Println(ok, er)

	result, msg, err := checker.SuperCheck(user)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("-------SuperCheck(in interface{})---------------------------------------------------------------------------------------------")

	fmt.Println("匹配结果:", result, "信息:", msg)

	fmt.Println("-------FormatCheck(in interface{})/Validate(in interface{})---------------------------------------------------------------------------------------------")
	ok, msg, er = checker.FormatCheck(user)
	if er != nil {
		log.Println(er.Error())
	}
	fmt.Println("格式验证结果:", ok, "msg:", msg)

	fmt.Println("-------listAll---------------------------------------------------------------------------------------------")
	checker.ListAll()
}
