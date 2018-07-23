package ndr

import (
	"reflect"
	"testing"
)

type TmpStruct struct {
	A [][][][]int
}

func TestTmp(t *testing.T) {
	var a TmpStruct
	l := []int{2, 3, 3, 2}
	v := reflect.ValueOf(&a.A)
	v = v.Elem()
	ty := v.Type()
	t.Logf("%+v\n", ty)

	s := reflect.MakeSlice(ty, l[0], l[0])

	v.Set(s)
	makeSubSlices(v, l[1:])
	t.Logf("%+v\n", a)

	//a := [][][][]int{}
	//l := []int{2,3,3,2}
	//ty := reflect.ValueOf(a).Type()
	//t.Logf("Slice type: %+v\n", ty)
	//b := reflect.MakeSlice(ty, l[0], l[0])
	//t.Logf("Slice instance: %+v\n", b)
	//for i:=1; i<4; i++ {
	//	//b := reflect.MakeSlice(ty, l[0], l[0])
	//	ty = ty.Elem()
	//	t.Logf("%v %v\n", ty, ty.Kind())
	//	for j := 0 ; j < l[i-1]; j++ {
	//		b.Index(j).Set(reflect.MakeSlice(ty, l[i], l[i]))
	//		t.Logf("d %d i %d len %d\n", i, j, b.Index(j).Len())
	//	}
	//	t.Logf("b %+v\n", b)
	//}
	//t.Logf("%+v\n", b)
	//
	//tp := reflect.ValueOf(a).Type()
	//for i := 0; i < 4; i++ {
	//	tp = reflect.SliceOf(tp)
	//	t.Logf("%v %v\n", tp, tp.Kind())
	//}

	//a[0][0][0] = []int{1}
	//a[0][0][1] = 2
	//a[0][1][0] = 11
	//t.Logf("%v\n", reflect.ValueOf(a).Kind())
	//
	//t.Logf("%v\n", reflect.TypeOf(a).Elem().Kind())
	//
	//t.Logf("%v\n", reflect.ValueOf(a).Type().Elem().Elem())        //new v
	//
	//
	//t.Logf("%v\n", reflect.ValueOf(a).Index(0))        //new v
	//t.Logf("%v\n", reflect.ValueOf(a).Index(0).Type()) //new t
	//
	//i := []int{2, 3, 2}
	//t.Logf("%v\n", len(multiDimensionalIndexPermutations(i[:len(i)-1])))
	//var s string
	//for i := 1; i <= 12; i++ {
	//	s = fmt.Sprintf("%s%s000000", s, hex.EncodeToString([]byte{byte(i)}))
	//}
	//t.Logf("%s\n", s)
}
