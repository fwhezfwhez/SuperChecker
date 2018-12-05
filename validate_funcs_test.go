package superChecker

import (
	"fmt"
	"testing"
)

func TestChineseOnly(t *testing.T) {
	type User struct {
		ChineseName string `validate:"func,chineseOnly"`
	}
	user := User{ChineseName: "中文"}
	checker := GetChecker()
	checker.AddFunc(ChineseOnly, "chineseOnly")
	ok, msg, er := checker.Validate(user)
	if er != nil {
		panic(er)
	}
	if !ok {
		fmt.Println(msg)
		return
	}
	fmt.Println("success")
}
