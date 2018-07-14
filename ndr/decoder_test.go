package ndr

import (
	"bytes"
	"encoding/hex"
	"testing"

	"gopkg.in/jcmturner/gokrb5.v5/mstypes"
)

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

func TestDecode(t *testing.T) {
	hexstr := "01100800cccccccca00400000000000000000200d186660f"
	b, _ := hex.DecodeString(hexstr)
	ft := new(mstypes.FileTime)
	dec := NewDecoder(bytes.NewReader(b), 4)
	err := dec.Decode(ft)
	if err != nil {
		t.Fatalf("error decoding: %v", err)
	}
	t.Logf("FT: %+v %v\n", ft, ft.Time())

	p := 20
	fto := mstypes.ReadFileTime(&b, &p, &dec.ch.Endianness)
	t.Logf("FTO: %+v %v\n", fto, fto.Time())
}
