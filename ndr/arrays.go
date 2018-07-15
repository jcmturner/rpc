package ndr

import (
	"fmt"
	"reflect"
)

// parseDimensions returns the number of dimensions, the size of each dimension and type of the member at the deepest level.
func parseDimensions(v reflect.Value, t reflect.Type) (n int, l []int, tb reflect.Type) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return
	}
	n++
	l = append(l, v.Len())
	if t.Elem().Kind() == reflect.Array || t.Elem().Kind() == reflect.Slice {
		// contains array or slice
		var s int
		var m []int
		s, m, tb = parseDimensions(v.Index(0), t.Elem())
		l = append(l, m...)
		n += s
	} else {
		tb = t.Elem()
	}
	return
}

// readUniDimensionalFixedArray reads an array (not slice) from the byte stream.
func (dec *Decoder) fillUniDimensionalFixedArray(v reflect.Value, tag reflect.StructTag) (err error) {
	//vptr := reflect.ValueOf(a)
	//if vptr.Kind() != reflect.Ptr || vptr.Elem().Kind() != reflect.Array  {
	//	err = errors.New("cannot fill uni-dimensional fixed array, a pointer to a Go array must be provided")
	//}
	//v := vptr.Elem()  // value == array kind
	// fill each index
	for i := 0; i < v.Len(); i++ {
		err = dec.fill(v.Index(i), tag)
		if err != nil {
			err = fmt.Errorf("could not fill index %d of uni-dimensional fixed array: %v", i, err)
			return
		}
	}
	return
}
