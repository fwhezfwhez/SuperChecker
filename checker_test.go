package superChecker

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"log"
	"strconv"
	"testing"
	"time"
)

// explain flag's usages
func TestChecker_ValidateOne(t *testing.T) {
	sp := GetChecker()

	var ok bool
	var msg string
	var e error

	log.SetFlags(log.Lshortfile)
	// validate flag regex
	ok, msg, e = sp.ValidateOne("Username", "fetdsfd", "regex,^[\u4E00-\u9FA5da-zA-Z]*$")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 1", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Username", "fetds####fd", "regex,^[\u4E00-\u9FA5da-zA-Z]*$")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 2", ok, msg, e)
	}

	// validate flag range,in
	ok, msg, e = sp.ValidateOne("Username", "admin", "range,[admin,administrator,root]")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 3", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Username", "admindd", "range,[admin,administrator,root]")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 4", ok, msg, e)
	}
	// validate flag int,float
	ok, msg, e = sp.ValidateOne("Age", 4, "int,3:5")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 5", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 4, "int,3:")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 6", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 4, "int,:5")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 7", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Age", 4.1, "float,3:5")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 8", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 4.1, "float,3:")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 9", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Age", 4.1, "float,:5")
	log.Println(ok, msg, e)

	if !ok || e != nil {
		t.Fatal("spot 10", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 1, "int,3:5")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 11", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 2, "int,3:")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 12", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 6, "int,:5")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 13", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Age", 5.2, "float,3:5")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 14", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 1.7, "float,3:")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 15", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Age", 5.9, "float,:5")
	log.Println(ok, msg, e)

	if ok {
		t.Fatal("spot 16", ok, msg, e)
	}

	// validate flag func, function
	if e := sp.AddFunc(func(in interface{}, fieldName string) (bool, string, error) {
		age := in.(int)

		if age > 0 && age < 200 {
			return true, "success", nil
		}
		return false, "age require between 0 and 200 but got " + strconv.Itoa(age), nil
	}, "lengthBetween0And200"); e != nil {
		t.Fatal("spot 17: ", ok, msg, e)
	}
	ok, msg, e = sp.ValidateOne("Age", 20, "func,lengthBetween0And200")
	log.Println(ok, msg, e)
	if !ok || e != nil {
		t.Fatal("spot 17", ok, msg, e)
	}

	ok, msg, e = sp.ValidateOne("Age", 2001, "func,lengthBetween0And200")
	log.Println(ok, msg, e)
	if ok {
		t.Fatal("spot 18", ok, msg, e)
	}
}

func TestChecker_ValidateByTagKeyAndMapValue(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	type User struct {
		Age       *int32    `json:"age,omitempty"`
		Username  *string   `json:"username,omitempty"`
		CreatedAt *string   `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Salary    *float64  `json:"salary,omitempty"`
		State     *int32    `json:"state,omitempty"`
	}
	user := User{
		Username:  proto.String("superchecker"),
		Age:       proto.Int32(20),
		CreatedAt: proto.String("2005-01-01"),
		UpdatedAt: time.Now(),
		Salary:    proto.Float64(2000),
		State:     proto.Int32(3),
	}

	sp := GetChecker()
	ok, msg, e := sp.ValidateByTagKeyAndMapValue(user, "json", map[string]string{
		"username":   "regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$",
		"age":        "int,0:200",
		"created_at": "time.Time,2006-01-02",
		"updated_at": "time.Time,2006-01-02",
		"salary":     "float,0:",
		"state":      "range,[1,2,3]",
	})
	log.Println(ok, msg, e)
	if !ok || e != nil {
		t.Fatal("spot 1", ok, msg, e)
	}

	//user.Age = proto.Int32(3000)
	//user.Username = proto.String("#$$#!")
	user.State = proto.Int32(5)
	ok, msg, e = sp.ValidateByTagKeyAndMapValue(user, "json", map[string]string{
		"username":   "regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$",
		"age":        "int,0:200",
		"created_at": "time.Time,2006-01-02",
		"updated_at": "time.Time,2006-01-02",
		"salary":     "float,0:",
		"state":      "range,[1,2,3]",
	})
	log.Println(ok, msg, e)
	if ok {
		t.Fatal("spot 2", ok, msg, e)
	}

}

type Week int

const (
	Monday    Week = 0
	Tuesday   Week = 1
	Wednsday  Week = 2
	Thursday  Week = 3
	Friday    Week = 4
	Satuarday Week = 5
	Sunday    Week = 6
)

type Employee struct {
	RestDay *Week `json:"rest_day"`
}

func (e *Employee) String() string {
	return "employee"
}
func TestEnum(t *testing.T) {
	var day = Friday
	var em = Employee{
		RestDay: &day,
	}
	sp := GetChecker()
	ok, msg, e := sp.ValidateByTagKeyAndMapValue(em, "json", map[string]string{
		"rest_day": "range,[1:7]",
	})
	fmt.Println(ok, msg, e)
}

func TestRange(t *testing.T){
	sp := GetChecker()
	type ShopConfig struct {
		Id            int             `gorm:"column:id;default:" json:"id" form:"id"`
		GameId        int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		PlatformId    int             `gorm:"column:platform_id;default:" json:"platform_id" form:"platform_id"`
		PageType      int             `gorm:"column:page_type;default:" json:"page_type" form:"page_type"`
		PropConfigId  int             `gorm:"column:prop_config_id;default:" json:"prop_config_id" form:"prop_config_id"`
		PropPriceMeal json.RawMessage `gorm:"column:prop_price_meal;default:" json:"prop_price_meal" form:"prop_price_meal"`
		Position      int             `gorm:"column:position;default:" json:"position" form:"position"`
		Status        int             `gorm:"column:status;default:" json:"status" form:"status"`
	}
	sc := ShopConfig{
		PropPriceMeal:json.RawMessage(`{
         "cost_num": 10
         }`),
	}
	type PriceMeal struct {
		Type       string `json:"type"`
		CostNum    int    `json:"cost_num"`
		PropBuyNum int    `json:"prop_buy_num"`
		PropBuyDay int    `json:"prop_buy_day"`
	}
	var pm PriceMeal
	e := json.Unmarshal(sc.PropPriceMeal, &pm)
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
		return
	}
	//ok,msg,e := sp.ValidateByTagKeyAndMapValue(pm, "json", map[string]string{
	//	"type":"range,[gold,money,diamond]",
	//	"cost_num":"int,0:",
	//	"prop_buy_num":"int,0:",
	//	"prop_buy_day":"int,-1:",
	//})
	ok,msg,e := sp.ValidateOne("type", pm.Type, "range,[gold,money,diamond]")
	if e!=nil {
		fmt.Println(e.Error())
		t.Fail()
		return
	}
	if !ok {
		fmt.Println(msg)
		t.Fail()
		return
	}

}

func TestRangeZero(t *testing.T) {
	sp := GetChecker()
	ok, msg, e := sp.ValidateOne("type", 5, "range,[0]")
	fmt.Println(ok,msg,e)

	ok, msg, e = sp.ValidateOne("type", "", "range,[]")
	fmt.Println(ok,msg,e)
}
