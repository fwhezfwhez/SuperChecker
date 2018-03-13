一款轻巧方便的注解式数据验证
##Example
```go
    package main
    
    import (
        "superChecker"
        "fmt"
        "github.com/lunny/log"
    )
    //superChecker的标签值就是需要匹配的key，逗号","表示与,"|"表示或
    type User struct {
        UserName string `superChecker:"userName" json:"userName"`
        Password string `superChecker:"password"`
        Phone string `superChecker:"mobilephone|telephone"`
        Text string `superChecker:"length,chineseOnly,notNull"`
    }
    func main(){
        user := User{
            "fwhez",
            "a1dfdasfsdf",
            "12578854875",
            "中",
        }
        //获取checker对象
        checker :=superChecker.GetChecker()
        
        //添加自定义正则,key值大小写不敏感，与tag的标签值对应
        //细心的读者可以发现，不存在key为username的regex，为什么可以匹配上呢?
        //因为checker有两个匹配池，默认添加的都是自定义池，默认池在init.go里，可以去看看
        //匹配优先选取自定义池，也就是说，当自定义池和默认池同时定义username，则会取自定义池的进行匹配
        checker.AddRegex("passWoRd","^[\\s\\S]{6,}$")
        checker.AddRegex("length","^[\\s\\S]{0,20}$")
        checker.AddRegex("chineseOnly","^[\u4E00-\u9FA5]*$")
        result,msg,err :=checker.SuperCheck(user)
        if err!=nil {
            log.Println(err)
        }
        fmt.Println("匹配结果:",result,"信息:",msg)
    }
```

####使用步骤
1. 给需要验证的结构体添加superChecker注解与Tag值，类比于userName的json注解,
值表示正则表达式的索引key
```go
   type User struct {
   	UserName string `superChecker:"userName"`
   	Password string `superChecker:"password"`
   	Phone string `superChecker:"mobilephone|telephone"`
   	Text string `superChecker:"length,chineseOnly,notNull"`
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
