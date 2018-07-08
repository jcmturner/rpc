package ndr

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

type NDRStruct struct {
	Count    int64 `ndr:"value,key:val"`
	Num      uint32
	ArrayPtr Ptr
}

func TestTemp(t *testing.T) {
	var n NDRStruct
	//dec := NewDecoder(bytes.NewReader([]byte{1,0,0,0,20,0,1,0}))

	t.Logf("Type: %+v\n", reflect.TypeOf(n))
	t.Logf("Kind %+v\n", reflect.TypeOf(n).Kind())
	t.Logf("Name %+v\n", reflect.TypeOf(n).Name())
	t.Logf("Field 0 %+v\n", reflect.TypeOf(n).Field(0))
	t.Logf("Field 1 %+v\n", reflect.TypeOf(n).Field(1))
	t.Logf("NumField %+v\n", reflect.TypeOf(n).NumField())

	t.Logf("Value: %+v\n", reflect.ValueOf(&n).Elem())
	reflect.ValueOf(&n).Elem().Field(1).Set(reflect.ValueOf(uint32(1)))

	t.Logf("ST: %+v\n", n)
}

func TestReadCommonHeader(t *testing.T) {
	var tests = []struct {
		EncodedHex string
		ExpectFail bool
	}{
		{"01100800cccccccc", false}, // Little Endian
		{"01000008cccccccc", false}, // Big Endian have to change the bytes for the header size? This test vector was artificially created. Need proper test vector
		//{"01100800cccccccc1802000000000000", false},
		//{"01100800cccccccc0002000000000000", false},
		//{"01100800cccccccc0001000000000000", false},
		//{"01100800cccccccce000000000000000", false},
		//{"01100800ccccccccf000000000000000", false},
		//{"01100800cccccccc7801000000000000", false},
		//{"01100800cccccccc4801000000000000", false},
		//{"01100800ccccccccd001000000000000", false},
		{"02100800cccccccc", true}, // Incorrect version
		{"02100900cccccccc", true}, // Incorrect length

	}

	for i, test := range tests {
		b, _ := hex.DecodeString(test.EncodedHex)
		dec := NewDecoder(bytes.NewReader(b), 4)
		err := dec.readCommonHeader()
		if err != nil && !test.ExpectFail {
			t.Errorf("error reading common header of test %d: %v", i, err)
		}
		if err == nil && test.ExpectFail {
			t.Errorf("expected failure on reading common header of test %d: %v", i, err)
		}
	}
}

func TestReadPrivateHeader(t *testing.T) {
	var tests = []struct {
		EncodedHex string
		ExpectFail bool
		Length     int
	}{
		{"01100800cccccccc1802000000000000", false, 536},
		{"01100800cccccccc0002000000000000", false, 512},
		{"01100800cccccccc0001000000000000", false, 256},
		{"01100800ccccccccFF00000000000000", true, 255}, // Length not multiple of 8
		{"01100800cccccccc00010000000000", true, 256},   // Too short

	}

	for i, test := range tests {
		b, _ := hex.DecodeString(test.EncodedHex)
		dec := NewDecoder(bytes.NewReader(b), 4)
		err := dec.readCommonHeader()
		if err != nil {
			t.Errorf("error reading common header of test %d: %v", i, err)
		}
		err = dec.readPrivateHeader()
		if err != nil && !test.ExpectFail {
			t.Errorf("error reading private header of test %d: %v", i, err)
		}
		if err == nil && test.ExpectFail {
			t.Errorf("expected failure on reading private header of test %d: %v", i, err)
		}
		if dec.ph.ObjectBufferLength != uint32(test.Length) {
			t.Errorf("Objectbuffer length expected %d actual %d", test.Length, dec.ph.ObjectBufferLength)
		}
	}
}
