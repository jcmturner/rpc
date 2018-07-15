package ndr

import (
	"errors"
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
func (dec *Decoder) fillUniDimensionalFixedArray(v reflect.Value, tag reflect.StructTag) error {
	for i := 0; i < v.Len(); i++ {
		err := dec.fill(v.Index(i), tag)
		if err != nil {
			return fmt.Errorf("could not fill index %d of uni-dimensional fixed array: %v", i, err)
		}
	}
	return nil
}

func (dec *Decoder) fillUniDimensionalConformantArray(v reflect.Value, tag reflect.StructTag) error {
	s, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not establish size of uni-dimensional conformant array: %v", err)
	}
	n := int(s)
	a := reflect.MakeSlice(v.Type(), n, n)
	for i := 0; i < n; i++ {
		err := dec.fill(a.Index(i), tag)
		if err != nil {
			return fmt.Errorf("could not fill index %d of uni-dimensional conformant array: %v", i, err)
		}
	}
	v.Set(a)
	return nil
}

func (dec *Decoder) fillUniDimensionalVaryingArray(v reflect.Value, tag reflect.StructTag) error {
	o, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not read offset of uni-dimensional varying array: %v", err)
	}
	s, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not establish size of uni-dimensional varying array: %v", err)
	}
	t := v.Type()
	os := t.Elem().Size()
	_, err = dec.r.Discard(int(o) * int(os))
	if err != nil {
		return fmt.Errorf("could not shift offset to read uni-dimensional varying array: %v", err)
	}
	n := int(s)
	a := reflect.MakeSlice(t, n, n)
	for i := 0; i < n; i++ {
		err := dec.fill(a.Index(i), tag)
		if err != nil {
			return fmt.Errorf("could not fill index %d of uni-dimensional varying array: %v", i, err)
		}
	}
	v.Set(a)
	return nil
}

func (dec *Decoder) fillUniDimensionalConformantVaryingArray(v reflect.Value, tag reflect.StructTag) error {
	m, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not read max count of uni-dimensional varying array: %v", err)
	}
	o, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not read offset of uni-dimensional varying array: %v", err)
	}
	s, err := dec.readUint32()
	if err != nil {
		return fmt.Errorf("could not establish size of uni-dimensional varying array: %v", err)
	}
	if m < o+s {
		return errors.New("max count is less than the offset plus actual count")
	}
	t := v.Type()
	os := t.Elem().Size()
	_, err = dec.r.Discard(int(o) * int(os))
	if err != nil {
		return fmt.Errorf("could not shift offset to read uni-dimensional varying array: %v", err)
	}
	n := int(s)
	a := reflect.MakeSlice(t, n, n)
	for i := 0; i < n; i++ {
		err := dec.fill(a.Index(i), tag)
		if err != nil {
			return fmt.Errorf("could not fill index %d of uni-dimensional varying array: %v", i, err)
		}
	}
	v.Set(a)
	return nil
}
