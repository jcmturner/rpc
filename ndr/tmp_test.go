package ndr

import (
	"reflect"
	"testing"
)

func TestTmp(t *testing.T) {
	a := [2][2]int{}
	t.Logf("%v\n", reflect.ValueOf(a).Kind())

	t.Logf("%v\n", reflect.TypeOf(a).Elem().Kind())

	t.Logf("%v\n", reflect.ValueOf(a[0]))

}
