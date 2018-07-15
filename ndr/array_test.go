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
