一款轻巧方便的注解式数据验证
##Example
```go
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
```

####使用步骤
1. 给需要验证的结构体添加superChecker注解与Tag值，类比于userName的json注解,
值表示正则表达式的索引key
**目前只支持int,float,time.Time的Validate,int和float取区间number1:number2(包括边界),时间直接书写格式，不限定具体时间**
```go
   type User struct {
	UserName string `superChecker:"userName" json:"userName" `
	Password string `superChecker:"password"`
	Phone string  `superChecker:"mobilephone|telephone"`
	Text string //`superChecker:"length,chineseOnly,notNull"`

	Age string `validate:"int,0:200"`
	Salary string `validate:"float,0:"`
	InTime string `validate:"time.Time,2006/1/2 15:04:05"`
   }
```
2. 创建checker对象，并添加以tag值为key的正则表达式
 ```go
    //获取checker对象
     checker :=superChecker.GetChecker()
     checker.AddRegex("passWoRd","^[\\s\\S]{6,}$")
     checker.AddRegex("length","^[\\s\\S]{0,20}$")
```
3.进行匹配,并输出结果
```go
     result,msg,err :=checker.SuperCheck(user)
            if err!=nil {
                log.Println(err)
            }
            fmt.Println("匹配结果:",result,"信息:",msg)
      ok,msg,er =checker.FormatCheck(user) //checker.Validate(user)  is the same
      if er!=nil{
            		fmt.Println(er.Error())
            		return
       }
       fmt.Println("格式验证结果:",ok,"msg:",msg)
```
结果格式:
失败
```go
分配成功
注入默认成功
匹配结果: false 信息: Phone 匹配失败
```
成功
```
分配成功
注入默认成功
匹配结果: true 信息: 匹配成功
```

##特殊操作
```go
        //1.给默认池添加正则
        checker.AddDefaultRegex("chineseOnly","^[\u4E00-\u9FA5]*$")
    
        //2.显示所有默认池k，v
        checker.ListDefault()
        //3.显示所有自定义池k,v
        checker.ListRegexBuilder()
        //4.显示上面两个
        checker.ListAll()
```
**注意**：
1. checker里存放的正则是编译过的，而不是string格式,所以打印出来的v值没有意义，用户可以看看k有哪些提供参考
2. 可以同时存在多个","或者多个"|"，但是不能,|同时出现在一个标签里
3. 不支持括号运算级别
