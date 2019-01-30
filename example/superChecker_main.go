// Tips:
// 'superChecker' tag has been stopped developing, the old functions will be remained.
// 'validate' tag is to replace the old ussages:
// "Name string `superChecker:"key1,key2"`" equals to "Name string `validate:"regex,key1,key2"`"

package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"superChecker"
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