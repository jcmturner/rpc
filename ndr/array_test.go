package ndr

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StructWithArray struct {
	A [4]uint32
}

type StructWithConformantSlice struct {
	A []uint32 `ndr:"conformant"`
}

type StructWithVaryingSlice struct {
	A []uint32 `ndr:"varying"`
}

type StructWithConformantVaryingSlice struct {
	A []uint32 `ndr:"conformant,varying"`
}

func TestReadUniDimensionalFixedArrary(t *testing.T) {
	hexStr := "01100800cccccccca0040000000000000000020001000000020000000300000004000000"
	b, _ := hex.DecodeString(hexStr)
	a := new(StructWithArray)
	dec := NewDecoder(bytes.NewReader(b), 4)
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i := range a.A {
		assert.Equal(t, uint32(i+1), a.A[i], "Value of index %d not as expected", i)
	}
}

func TestReadUniDimensionalConformantArrary(t *testing.T) {
	hexStr := "01100800cccccccca004000000000000000002000400000001000000020000000300000004000000"
	b, _ := hex.DecodeString(hexStr)
	a := new(StructWithConformantSlice)
	dec := NewDecoder(bytes.NewReader(b), 4)
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i := range a.A {
		assert.Equal(t, uint32(i+1), a.A[i], "Value of index %d not as expected", i)
	}
}

func TestReadUniDimensionalVaryingArrary(t *testing.T) {
	hexStr := "01100800cccccccca004000000000000000002000200000004000000000000000000000001000000020000000300000004000000"
	b, _ := hex.DecodeString(hexStr)
	a := new(StructWithVaryingSlice)
	dec := NewDecoder(bytes.NewReader(b), 4)
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i := range a.A {
		assert.Equal(t, uint32(i+1), a.A[i], "Value of index %d not as expected", i)
	}
}

func TestReadUniDimensionalConformantVaryingArrary(t *testing.T) {
	hexStr := "01100800cccccccca00400000000000000000200060000000200000004000000000000000000000001000000020000000300000004000000"
	b, _ := hex.DecodeString(hexStr)
	a := new(StructWithConformantVaryingSlice)
	dec := NewDecoder(bytes.NewReader(b), 4)
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i := range a.A {
		assert.Equal(t, uint32(i+1), a.A[i], "Value of index %d not as expected", i)
	}
}
