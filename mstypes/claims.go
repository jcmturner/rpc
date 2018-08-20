package mstypes

import (
	"bytes"
	"errors"

	"github.com/jcmturner/rpc/ndr"
)

// Compression format assigned numbers.
const (
	CompressionFormatNone       uint16 = 0
	CompressionFormatLZNT1      uint16 = 2
	CompressionFormatXPress     uint16 = 3
	CompressionFormatXPressHuff uint16 = 4
)

// ClaimsSourceTypeAD https://msdn.microsoft.com/en-us/library/hh553809.aspx
const ClaimsSourceTypeAD uint16 = 1

// Claim Type assigned numbers
const (
	ClaimTypeIDInt64    uint16 = 1
	ClaimTypeIDUInt64   uint16 = 2
	ClaimTypeIDString   uint16 = 3
	ClaimsTypeIDBoolean uint16 = 6
)

// ClaimsBlob implements https://msdn.microsoft.com/en-us/library/hh554119.aspx
type ClaimsBlob struct {
	Size        uint32
	EncodedBlob EncodedBlob
}

type EncodedBlob []byte

func (b EncodedBlob) Size(c interface{}) int {
	cb := c.(ClaimsBlob)
	return int(cb.Size)
}

// ClaimsSetMetadata implements https://msdn.microsoft.com/en-us/library/hh554073.aspx
type ClaimsSetMetadata struct {
	claimsSetSize             uint32
	ClaimsSetBytes            ClaimsSetBytes `ndr:"pointer"`
	CompressionFormat         uint16         // Enum see constants for options
	uncompressedClaimsSetSize uint32
	ReservedType              uint16
	reservedFieldSize         uint32
	ReservedField             ClaimsSetMetadataReservedFieldBytes `ndr:"pointer"`
}

type ClaimsSetBytes []byte

func (b ClaimsSetBytes) Size(p interface{}) int {
	c := p.(ClaimsSetMetadata)
	return int(c.claimsSetSize)
}

type ClaimsSetMetadataReservedFieldBytes []byte

func (b ClaimsSetMetadataReservedFieldBytes) Size(p interface{}) int {
	c := p.(ClaimsSetMetadata)
	return int(c.reservedFieldSize)
}

func (m *ClaimsSetMetadata) ClaimsSet() (c ClaimsSet, err error) {
	if len(m.ClaimsSetBytes) < 1 {
		err = errors.New("no bytes available for ClaimsSet")
		return
	}
	// TODO switch statement to decompress ClaimsSetBytes
	if m.CompressionFormat != CompressionFormatNone {
		err = errors.New("compressed ClaimsSet not currently supported")
		return
	}
	dec := ndr.NewDecoder(bytes.NewReader(m.ClaimsSetBytes))
	err = dec.Decode(&c)
	return
}

// ClaimsSet implements https://msdn.microsoft.com/en-us/library/hh554122.aspx
type ClaimsSet struct {
	ClaimsArrayCount  uint32
	ClaimsArrays      []ClaimsArray `ndr:"pointer,conformant"`
	ReservedType      uint16
	reservedFieldSize uint32
	ReservedField     ClaimsSetReservedFieldBytes `ndr:"pointer"`
}

type ClaimsSetReservedFieldBytes []byte

func (b ClaimsSetReservedFieldBytes) Size(p interface{}) int {
	c := p.(ClaimsSet)
	return int(c.reservedFieldSize)
}

// ClaimsArray implements https://msdn.microsoft.com/en-us/library/hh536458.aspx
type ClaimsArray struct {
	ClaimsSourceType uint16
	ClaimsCount      uint32
	ClaimEntries     []ClaimEntry `ndr:"pointer,encapsulated"`
}

// ClaimEntry is a NDR union that implements https://msdn.microsoft.com/en-us/library/hh536374.aspx
type ClaimEntry struct {
	ID         string
	Tag        uint16
	TypeInt64  ClaimTypeInt64
	TypeUInt64 ClaimTypeUInt64
	TypeString ClaimTypeString
	TypeBool   ClaimTypeBoolean
}

func (u ClaimEntry) SwitchFunc(tag interface{}) string {
	t := tag.(uint16)
	switch t {
	case ClaimTypeIDInt64:
		return "TypeInt64"
	case ClaimTypeIDUInt64:
		return "TypeUInt64"
	case ClaimTypeIDString:
		return "TypeString"
	case ClaimsTypeIDBoolean:
		return "TypeBool"
	}
	return ""
}

// ClaimTypeInt64 is a claim of type int64
type ClaimTypeInt64 struct {
	ValueCount uint32
	Value      []int64 `ndr:"pointer,conformant"`
}

// ClaimTypeUInt64 is a claim of type uint64
type ClaimTypeUInt64 struct {
	ValueCount uint32
	Value      []uint64 `ndr:"pointer,conformant"`
}

// ClaimTypeString is a claim of type string
type ClaimTypeString struct {
	ValueCount uint32
	Value      []string `ndr:"pointer,conformant"`
}

// ClaimTypeBoolean is a claim of type bool
type ClaimTypeBoolean struct {
	ValueCount uint32
	Value      []bool `ndr:"pointer,conformant"`
}
