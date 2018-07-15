package ndr

import (
	"reflect"
	"testing"
)

func TestTmp(t *testing.T) {
	a := [2][2][2][]SimpleTest{}

	//a[0][0][0] = []int{1}
	//a[0][0][1] = 2
	//a[0][1][0] = 11
	t.Logf("%v\n", reflect.ValueOf(a).Kind())

	t.Logf("%v\n", reflect.TypeOf(a).Elem().Kind())

	t.Logf("%v\n", reflect.ValueOf(a).Index(0))        //new v
	t.Logf("%v\n", reflect.ValueOf(a).Index(0).Type()) //new t

	n, v, ta := parseDimensions(reflect.ValueOf(a), reflect.TypeOf(a))
	t.Logf("d: %v %v %v\n", n, v, ta)

}
