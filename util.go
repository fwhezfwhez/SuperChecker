package superChecker

import (
	"errors"
	"fmt"
	"github.com/fwhezfwhez/jsoncrack"
	"go/types"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// SmartPrint pretty format print an input value,which should be a struct
func SmartPrint(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("获取到数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}

// ToString Change arg to string
func ToString(arg interface{}, timeFormat ...string) string {
	if len(timeFormat) > 1 {
		log.SetFlags(log.Llongfile | log.LstdFlags)
		log.Println(errors.New(fmt.Sprintf("timeFormat's length should be one")))
	}
	switch v := arg.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeFormat) == 1 {
			return v.Format(timeFormat[0])
		}
		return v.Format("2006-01-02 15:04:05")
	case jsoncrack.Time:
		if len(timeFormat)==1 {
			return v.Time().Format(timeFormat[0])
		}
	    return v.Time().Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return v.String()
	case types.Pointer:
		return "not for ptr,you might need &ptr"
	default:
		return ""
	}
}

func RemovePrefix(s string, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return s
	}
	return s[len(prefix):]
}
