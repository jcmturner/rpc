package ndr

import (
	"encoding/hex"
	"fmt"
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

	i := []int{2, 3, 2}
	t.Logf("%v\n", len(multiDimensionalIndexPermutations(i[:len(i)-1])))
	var s string
	for i := 1; i <= 12; i++ {
		s = fmt.Sprintf("%s%s000000", s, hex.EncodeToString([]byte{byte(i)}))
	}
	t.Logf("%s\n", s)
}
