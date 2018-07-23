package ndr

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func intFromTag(tag reflect.StructTag, key string) (int, error) {
	ndrTag := parseTags(tag)
	d := 1
	if n, ok := ndrTag.Map[key]; ok {
		i, err := strconv.Atoi(n)
		if err != nil {
			return d, fmt.Errorf("invalid dimensions tag [%s]: %v", n, err)
		}
		d = i
	}
	return d, nil
}

// parseDimensions returns the a slice of the size of each dimension and type of the member at the deepest level.
func parseDimensions(v reflect.Value) (l []int, tb reflect.Type) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return
	}
	l = append(l, v.Len())
	if t.Elem().Kind() == reflect.Array || t.Elem().Kind() == reflect.Slice {
		// contains array or slice
		var m []int
		m, tb = parseDimensions(v.Index(0))
		l = append(l, m...)
	} else {
		tb = t.Elem()
	}
	return
}

func sliceDimensions(t reflect.Type) (d int) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		d++
		n := sliceDimensions(t.Elem())
		d += n
	}
	return d
}

// makeSubSlices is a deep recursive creation/initialisation of multi-dimensional slices.
// Takes the reflect.Value of the 1st dimension and a slice of the lengths of the sub dimensions
func makeSubSlices(v reflect.Value, l []int) {
	ty := v.Type().Elem()
	if ty.Kind() != reflect.Slice {
		return
	}
	for i := 0; i < v.Len(); i++ {
		s := reflect.MakeSlice(ty, l[0], l[0])
		v.Index(i).Set(s)
		// Are there more sub dimensions?
		if len(l) > 1 {
			makeSubSlices(v.Index(i), l[1:])
		}
	}
	return
}

// multiDimensionalIndexPermutations returns all the permutations of the indexes of a multi-dimensional slice.
// The input is a slice of integers that indicates the max size/length of each dimension
func multiDimensionalIndexPermutations(l []int) (ps [][]int) {
	z := make([]int, len(l), len(l)) // The zeros permutation
	ps = append(ps, z)
	// for each dimension, in reverse
	for i := len(l) - 1; i >= 0; i-- {
		ws := make([][]int, len(ps))
		copy(ws, ps)
		//create a permutation for each of the iterations of the current dimension
		for j := 1; j <= l[i]-1; j++ {
			// For each existing permutation
			for _, p := range ws {
				np := make([]int, len(p), len(p))
				copy(np, p)
				np[i] = j
				ps = append(ps, np)
			}
		}
	}
	return
}

// readUniDimensionalFixedArray reads an array (not slice) from the byte stream.
func (dec *Decoder) fillUniDimensionalFixedArray(v reflect.Value, tag reflect.StructTag) error {
	for i := 0; i < v.Len(); i++ {
		err := dec.fill(v.Index(i), tag)
		if err != nil {
			return fmt.Errorf("could not fill index %d of fixed array: %v", i, err)
		}
	}
	return nil
}

func (dec *Decoder) fillFixedArray(v reflect.Value, tag reflect.StructTag) error {
	l, _ := parseDimensions(v)
	if len(l) < 1 {
		return errors.New("could not establish dimensions of fixed array")
	}
	if len(l) == 1 {
		err := dec.fillUniDimensionalFixedArray(v, tag)
		if err != nil {
			return fmt.Errorf("could not fill uni-dimensional fixed array: %v", err)
		}
		return nil
	}
	ps := multiDimensionalIndexPermutations(l[:len(l)-1])
	for _, p := range ps {
		// Get current multi-dimensional index to fill
		a := v
		for _, i := range p {
			a = a.Index(i)
		}
		// fill with the last dimension array
		err := dec.fillUniDimensionalFixedArray(a, tag)
		if err != nil {
			return fmt.Errorf("could not fill dimension %v of multi-dimensional fixed array: %v", p, err)
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

func (dec *Decoder) fillMultiDimensionalConformantArray(v reflect.Value, d int, tag reflect.StructTag) error {
	// Read the size of each dimensions from the ndr stream
	l := make([]int, d, d)
	for i := range l {
		s, err := dec.readUint32()
		if err != nil {
			return fmt.Errorf("could not read size of dimension %d: %v", i+1, err)
		}
		l[i] = int(s)
	}

	// Initialise size of slices
	// Initialise the size of the 1st dimension
	ty := v.Type()
	v.Set(reflect.MakeSlice(ty, l[0], l[0]))
	// Initialise the size of the other dimensions recursively
	makeSubSlices(v, l[1:])

	// Get all permutations of the indexes and go through each and fill
	ps := multiDimensionalIndexPermutations(l)
	for _, p := range ps {
		// Get current multi-dimensional index to fill
		a := v
		for _, i := range p {
			a = a.Index(i)
		}
		err := dec.fill(a, tag)
		if err != nil {
			return fmt.Errorf("could not fill index %v of multi-dimensional conformant array: %v", p, err)
		}
	}
	return nil
}

func (dec *Decoder) fillConformantArray(v reflect.Value, tag reflect.StructTag) error {
	d := sliceDimensions(v.Type())
	if d > 1 {
		err := dec.fillMultiDimensionalConformantArray(v, d, tag)
		if err != nil {
			return err
		}
	} else {
		err := dec.fillUniDimensionalConformantArray(v, tag)
		if err != nil {
			return err
		}
	}
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
