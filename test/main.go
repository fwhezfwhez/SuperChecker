package main

import (
	"superChecker"
	"fmt"
)

type User struct {
	UserName string `superChecker:"userName" json:"userName"`
	Password string `superChecker:"password"`
	Phone string `superChecker:"mobilephone|telephone"`
	Text string `superChecker:"length,chineseOnly,notNull"`
}
func main(){
	user := User{
		"fwhez",
		"a123gfdsd",
		"88545758",
		"你好",
	}
	checker :=superChecker.GetChecker()
	checker.AddRegex("passWoRd","^[\\s\\S]{6,}$")
	checker.AddRegex("length","^[\\s\\S]{0,20}$")
	checker.AddRegex("chineseOnly","^[\u4E00-\u9FA5]*$")
	result,err :=checker.SuperCheck(user)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("匹配结果:",result)
}