package superChecker

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"log"
	"reflect"
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
	fmt.Println(day)
	fmt.Println(reflect.TypeOf(day))
	fmt.Println(reflect.ValueOf(day).Interface())
	var em = Employee{
		RestDay: &day,
	}
	sp := GetChecker()
	ok, msg, e := sp.ValidateByTagKeyAndMapValue(em, "json", map[string]string{
		"rest_day": "range,[1:7]",
	})
	fmt.Println(ok, msg, e)
}
