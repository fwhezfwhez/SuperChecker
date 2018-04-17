package main

import (
	"superChecker"
	"fmt"
	"log"
)

type User struct {
	UserName string //`superChecker:"userName" json:"userName"`
	Password string //`superChecker:"password"`
	Phone string //`superChecker:"mobilephone|telephone"`
	Text string //`superChecker:"length,chineseOnly,notNull"`
}
func main(){
	user := User{
		"",
		"a1dfdasfsdf",
		"12578854875",
		"中",
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
}