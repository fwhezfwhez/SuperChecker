package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"superChecker"
	"time"
)

type Order struct {
	// TIME
	CreatedAt time.Time `validate:"time.time"`
	UpdatedAt string    `validate:"time.time,2006/01/02 15:04:05"`

	// INT
	Count    int `validate:"int,0:200"`
	MaxCount int `validate:"int,:200"`
	MinCount int `validate:"int,10:"`
	Count2   int `validate:"int64,0:200"`

	// FLOAT
	RewardRate    float64         `validate:"float,0:0.4"`
	MaxRewardRate float64         `validate:"float,:0.4"`
	MinRewradRate float64         `validate:"float,0:"`
	RewardRate2   float64         `validate:"float64,0:0.4"`
	RewardRate3   decimal.Decimal `validate:"decimal,0:0.4"`

	// REGEX
	OrderUsername  string `validate:"regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$"`
	OrderUsername2 string `validate:"regex,username"`

	// RANGE,IN
	OrderStatus     int    `validate:"range,[1,2,3,4]"`
	OrderStatusName string `validate:"in,[unpaid,paid,closed]"`

	// FUNC, FUNCTION
	MailTypeCheckBox  string `validate:"func,inAndLength,lengthMoreThan3"`
	MailTypeCheckBox2 string `validate:"function,lengthLessThan3|inAndLength"`
}
//func (o Order) XXSVValidateSVBCreate()(bool,string,error){
//	return true,"xxsvcreate wrong",nil
//}
//func (o Order) XXValidate()(bool,string,error){
//	return true,"xxv wrong",nil
//}
//func (o Order) XXSVValidate()(bool,string,error){
//	return true,"xxsv wrong",nil
//}

//func (o Order) XXValidateSVBCreate()(bool,string,error){
//	return true,"xxcreate wrong",nil
//}
//
//
//
//func (o Order) XXValidateSVBCreateSVSUpdate()(bool,string,error){
//	return true,"success",nil
//}
//func (o Order) XXSVValidateSVBCreateSVSUpdate()(bool,string,error){
//	return true,"xxsvcreateupdate wrong",nil
//}

func (o Order) ValidateSVBCreate()(bool,string,error){
	return false,"xxsvcreate wrong",nil
}

func main() {
	order := Order{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Format("2006/01/02 15:04:05"),

		Count:    200,
		MaxCount: 90,
		MinCount: 10,
		Count2:   100,

		RewardRate:    0.4,
		MaxRewardRate: 0.3,
		MinRewradRate: 0.1,
		RewardRate2:   0.1,
		RewardRate3:   decimal.NewFromFloat(0.4),

		OrderUsername:  "superCheckerValidate",
		OrderUsername2: "superCheckerValidate",

		OrderStatus:     3,
		OrderStatusName: "closed",

		MailTypeCheckBox:  "midMail",
		MailTypeCheckBox2: "midMail",
	}

	checker := superChecker.GetChecker()
	checker.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		v := superChecker.ToString(in)
		maxLength := 7
		if len(v) > maxLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' over length,want %d ,but got %d", fieldName, "inAndLength", maxLength, len(v)), nil
		}
		vrange := []string{"midMail", "shenMail", "yundaMail"}
		for _, value := range vrange {
			if value == v {
				return true, "success", nil
			}
		}
		return false, fmt.Sprintf("while validating field '%s', rule key '%s',  value '%s' not in '%v'", fieldName, "inAndLength", v, vrange), nil
	}, "inAndLength")
	checker.AddFunc(func(in interface{}, fieldName string)(bool, string, error){
		v := superChecker.ToString(in)
		minLength := 3
		if len(v) < minLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' too short length,want %d ,but got %d", fieldName, "inAndLength", minLength, len(v)), nil
		}
		return true, "success", nil
	},"lengthmorethan3")

	checker.AddFunc(func(in interface{}, fieldName string)(bool, string, error){
		v := superChecker.ToString(in)
		maxLength := 3
		if len(v) > maxLength {
			return false, fmt.Sprintf("while validating field '%s', rule key '%s' too short length,want %d ,but got %d", fieldName, "inAndLength", maxLength, len(v)), nil
		}
		return true, "success", nil
	},"lengthlessthan3")

	ok, msg, er := checker.Validate(order)
	if er != nil {
		fmt.Println(fmt.Sprintf("got an error, '%s'", er.Error()))
		return
	}
	if !ok {
		fmt.Println(fmt.Sprintf("validate fail because of '%s'", msg))
		return
	}


	// ioc, inverse of control
	// validate to combine as receiver to the dest struct
	ok, msg, er = checker.ValidateMethods(order,"create")
	if er != nil {
		fmt.Println(fmt.Sprintf("got an error, '%s'", er.Error()))
		return
	}
	if !ok {
		fmt.Println(fmt.Sprintf("validate fail because of '%s'", msg))
		return
	}
	fmt.Println("success")


}
