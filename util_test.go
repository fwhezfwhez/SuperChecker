package superChecker

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	var a = "1"
	fmt.Println(ToString(reflect.ValueOf(a)))
}
