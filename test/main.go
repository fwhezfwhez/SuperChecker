package main

import (
	"superChecker"
	"fmt"
	"log"
)

type User struct {
	UserName string `superChecker:"userName" json:"userName" `
	Password string `superChecker:"password"`
	Phone string  `superChecker:"mobilephone|telephone"`
	Text string //`superChecker:"length,chineseOnly,notNull"`

	Age string `validate:"int,0:200"`
	Salary string `validate:"float,0:"`
	InTime string `validate:"time.Time,2006/1/2 15:04:05"`
}
func main(){
	user := User{
		UserName:"d",
		Password:"a1dfdasfsdf",
		Phone:"undefine",
		Text:"undefined",
		Age:"200",
		Salary:"5",
		InTime:"2018/1/2 15:04:05",
	}
	checker :=superChecker.GetChecker()
	checker.AddRegex("passWoRd","^[\\s\\S]{6,}$")
	checker.AddRegex("length","^[\\s\\S]{0,20}$")
	checker.AddRegex("chineseOnly","^[\u4E00-\u9FA5]*$")
	result,msg,err :=checker.SuperCheck(user)
	if err!=nil {
		log.Println(err)
	}
	fmt.Println("匹配结果:",result,"信息:",msg)

	checker.AddDefaultRegex("chineseOnly","^[\u4E00-\u9FA5]*$")

	checker.ListDefault()

	checker.ListRegexBuilder()

	checker.ListAll()

	ok,er:=checker.Check("10000124","^[0-9]{8}$")
	fmt.Println(ok,er)

	ok,msg,er =checker.FormatCheck(user)
	if er!=nil{
		fmt.Println(er.Error())
		return
	}
	fmt.Println("格式验证结果:",ok,"msg:",msg)
}