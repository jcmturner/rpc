package ndr

import (
	"fmt"
)

// readConformantVaryingString reads a Conformant and Varying String from the bytes slice.
// A conformant and varying string is a string in which the maximum number of elements is not known beforehand and therefore is included in the representation of the string.
// NDR represents a conformant and varying string as an ordered sequence of representations of the string elements, preceded by three unsigned long integers.
// The first integer gives the maximum number of elements in the string, including the terminator.
// The second integer gives the offset from the first index of the string to the first index of the actual subset being passed.
// The third integer gives the actual number of elements being passed, including the terminator.
func (dec *Decoder) readConformantVaryingString() (string, error) {
	m, err := dec.readUint32() // Max element count
	if err != nil {
		return "", Malformed{fmt.Sprintf("could not read conformant varying string max count element: %v", err)}
	}
	o, err := dec.readUint32() // Offset
	if err != nil {
		return "", Malformed{fmt.Sprintf("could not read conformant varying string offset element: %v", err)}
	}
	a, err := dec.readUint32() // Actual count
	if err != nil {
		return "", Malformed{fmt.Sprintf("could not read conformant varying string actual count element: %v", err)}
	}
	if a > (m-o) || o > m {
		return "", Malformed{EText: fmt.Sprintf("not enough bytes to read conformant varying string. Max: %d, Offset: %d, Actual: %d", m, o, a)}
	}
	//Unicode string so each element is 2 bytes
	//move position based on the offset
	if o > 0 {
		_, err := dec.r.Discard(int(o * 2))
		if err != nil {
			return "", Malformed{"could not move to offset position to read conformant varying string"}
		}
	}
	s := make([]rune, a, a)
	for i := 0; i < len(s); i++ {
		r, err := dec.readUint16()
		if err != nil {
			return "", Malformed{fmt.Sprintf("could not read bytes for rune at index %d: %v", i, err)}
		}
		s[i] = rune(r)
	}
	dec.ensureAlignment()
	if len(s) > 0 {
		// Remove any null terminator
		if s[len(s)-1] == rune(0) {
			s = s[:len(s)-1]
		}
	}
	return string(s), nil
}

// readVaryingString reads a Conformant and Varying String from the bytes slice.
// NDR represents a varying string as an ordered sequence of representations of the string elements, preceded by two unsigned long integers.
// The first integer gives the offset from the first index of the string to the first index of the actual subset being
// passed. The second integer gives the actual number of elements being passed, including the terminator.
func (dec *Decoder) readVaryingString() (string, error) {
	o, err := dec.readUint32() // Offset
	if err != nil {
		return "", Malformed{fmt.Sprintf("could not read conformant varying string offset element: %v", err)}
	}
	a, err := dec.readUint32() // Actual count
	if err != nil {
		return "", Malformed{fmt.Sprintf("could not read conformant varying string actual count element: %v", err)}
	}
	//Unicode string so each element is 2 bytes
	//move position based on the offset
	if o > 0 {
		_, err := dec.r.Discard(int(o * 2))
		if err != nil {
			return "", Malformed{"could not move to offset position to read conformant varying string"}
		}
	}
	s := make([]rune, a, a)
	for i := 0; i < len(s); i++ {
		r, err := dec.readUint16()
		if err != nil {
			return "", Malformed{fmt.Sprintf("could not read bytes for rune at index %d: %v", i, err)}
		}
		s[i] = rune(r)
	}
	dec.ensureAlignment()
	if len(s) > 0 {
		// Remove any null terminator
		if s[len(s)-1] == rune(0) {
			s = s[:len(s)-1]
		}
	}
	return string(s), nil
}
