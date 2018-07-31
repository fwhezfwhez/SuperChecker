package main

import (
	"superChecker"
	"fmt"
	"log"
)

type User struct {
	InTime string `validate:"time.Time,2006/1/2 15:04:05"`
	K []int
	Phone string  `superChecker:"mobilephone|telephone"`
	UserName string `superChecker:"userName" json:"userName" `
	Password string `superChecker:"password"`

	Text string //`superChecker:"length,chineseOnly,notNull"`
	Number int `superChecker:"mobilephone"`

	Age string  `validate:"int"`
	Salary string `validate:"float"`

}
func main(){
	user := User{
		UserName:"d",
		Password:"a1dfdasfsdf",
		Phone:"13875847584",
		Text:"undefined",
		Age:"",
		Salary:"5",
		InTime:"2018-5-2",
		Number:9,
	}
	checker :=superChecker.GetChecker()
	checker.AddRegex("passWoRd","^[\\s\\S]{6,}$")
	checker.AddRegex("length","^[\\s\\S]{0,20}$")
	checker.AddRegex("chineseOnly","^[\u4E00-\u9FA5]*$")

	checker.AddDefaultRegex("chineseOnly","^[\u4E00-\u9FA5]*$")

	checker.ListDefault()

	checker.ListRegexBuilder()

	checker.ListAll()

	fmt.Println("-------Check(in string,rule string)---------------------------------------------------------------------------------------------")
	ok,er:=checker.Check("10000124","^[0-9]{8}$")
	fmt.Println(ok,er)

	result,msg,err :=checker.SuperCheck(user)
	if err!=nil {
		log.Println(err)
	}
	fmt.Println("-------SuperCheck(in interface{})---------------------------------------------------------------------------------------------")

	fmt.Println("匹配结果:",result,"信息:",msg)

	fmt.Println("-------FormatCheck(in interface{})---------------------------------------------------------------------------------------------")
	ok,msg,er =checker.FormatCheck(user)
	if er!=nil{
		log.Println(er.Error())
	}
	fmt.Println("格式验证结果:",ok,"msg:",msg)
}