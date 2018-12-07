package superChecker

import (
	"errors"
	"fmt"
	"strconv"
)

// only accept 'data' to be chinese
// use it like:
// type User struct{
//     ChineseName string `validate:"func,chineseOnly"`
// }
// func main(){
//     user := User{ChineseName: "ft"}
//     checker := superChecker.GetChecker()
//     checker.AddFunc(superChecker.ChineseOnly, "chineseOnly")
//     ok,msg,er :=checker.Validate(user)
//	   if er!=nil {
//        panic(er)
//     }
//     if !ok {
//         fmt.Println(msg)
//         return
//     }
//     fmt.Println("success")
// }
func ChineseOnly(data interface{}, fieldName string) (bool, string, error) {
	v := ToString(data)
	if data == "" {
		return true, "success", nil
	}
	if r, ok := compiledMap["chineseonly"]; !ok {
		panic(errors.New("chineseonly not found in compiledMap"))
	} else {
		if !r.MatchString(v) {
			return false, fmt.Sprintf("while validating field '%s',regex rule '%s',got unmatched value '%s'", fieldName, strconv.QuoteToASCII(r.String()), v), nil
		}
	}
	return true, v, nil
}


