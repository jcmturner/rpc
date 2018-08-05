package ndr

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestStr         = "hello world!"
	TestStrUTF16Hex = "680065006c006c006f00200077006f0072006c00640021000000"
)

type TestStructWithVaryingString struct {
	A string `ndr:"varying"`
}

type TestStructWithConformantVaryingString struct {
	A string `ndr:"conformant,varying"`
}

func Test_uint16SliceToString(t *testing.T) {
	b, _ := hex.DecodeString(TestStrUTF16Hex)
	var u []uint16
	for i := 0; i < len(b); i += 2 {
		u = append(u, binary.LittleEndian.Uint16(b[i:i+2]))
	}
	s := uint16SliceToString(u)
	assert.Equal(t, TestStr, s, "uint16SliceToString did not return as expected")
}

func Test_readVaryingString(t *testing.T) {
	ac := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(ac, uint32(len(TestStrUTF16Hex)/4))            // actual count of number of uint16 bytes
	hexStr := TestHeader + "00000000" + hex.EncodeToString(ac) + TestStrUTF16Hex // header:offset(0):actual count:data
	b, _ := hex.DecodeString(hexStr)
	a := new(TestStructWithVaryingString)
	dec := NewDecoder(bytes.NewReader(b), 1) //TODO why alignment = 4 fails
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal(t, TestStr, a.A, "value of decoded varying string not as expected")
}

func Test_readConformantVaryingString(t *testing.T) {
	ac := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(ac, uint32(len(TestStrUTF16Hex)/4))                                     // actual count of number of uint16 bytes
	hexStr := TestHeader + hex.EncodeToString(ac) + "00000000" + hex.EncodeToString(ac) + TestStrUTF16Hex // header:offset(0):actual count:data
	b, _ := hex.DecodeString(hexStr)
	a := new(TestStructWithConformantVaryingString)
	dec := NewDecoder(bytes.NewReader(b), 1) //TODO why alignment = 4 fails
	err := dec.Decode(a)
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal(t, TestStr, a.A, "value of decoded varying string not as expected")
}
