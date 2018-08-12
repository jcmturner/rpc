package ndr

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	TagConformant = "conformant"
	TagVarying    = "varying"
	TagPointer    = "pointer"
)

type Decoder struct {
	//mutex sync.Mutex    // each item must be received atomically
	r             *bufio.Reader // source of the data
	size          int           // initial size of bytes in buffer
	ch            CommonHeader  // NDR common header
	ph            PrivateHeader // NDR private header
	deferred      []deferedPtr
	conformantMax []uint32
	s             interface{} //pointer to the structure being populated
}

type deferedPtr struct {
	v   reflect.Value
	tag reflect.StructTag
}

func NewDecoder(r io.Reader) *Decoder {
	dec := new(Decoder)
	dec.r = bufio.NewReader(r)
	dec.r.Peek(int(commonHeaderBytes)) // For some reason an operation is needed on the buffer to initialise it so Buffered() != 0
	dec.size = dec.r.Buffered()
	return dec
}

func (dec *Decoder) Decode(s interface{}) error {
	dec.s = s
	err := dec.readCommonHeader()
	if err != nil {
		return err
	}
	err = dec.readPrivateHeader()
	if err != nil {
		return err
	}
	_, err = dec.r.Discard(4) //The next 4 bytes are an RPC unique pointer referent. We just skip these.
	if err != nil {
		return Malformed{fmt.Sprintf("unable to process byte stream: %v", err)}
	}
	// Scan for conformant fields as their max counts are moved to the beginning
	// http://pubs.opengroup.org/onlinepubs/9629399/chap14.htm#tagfcjh_37
	err = dec.conformantScan(s, reflect.StructTag(""))
	if err != nil {
		return fmt.Errorf("error scanning for conformant fields: %v", err)
	}
	for i := range dec.conformantMax {
		dec.conformantMax[i], err = dec.readUint32()
		if err != nil {
			return fmt.Errorf("could not read preceding conformant max count index %d: %v", i, err)
		}
	}
	err = dec.fill(s, reflect.StructTag(""), false)
	if err != nil {
		return err
	}
	// Read any deferred referents from pointers
	for _, p := range dec.deferred {
		err = dec.fill(p.v, p.tag, true)
		if err != nil {
			return fmt.Errorf("error filling deferred pointer referent: %v", err)
		}
	}
	return nil
}

func (dec *Decoder) conformantScan(s interface{}, tag reflect.StructTag) error {
	var v reflect.Value
	if r, ok := s.(reflect.Value); ok {
		v = r
	} else {
		if reflect.ValueOf(s).Kind() == reflect.Ptr {
			v = reflect.ValueOf(s).Elem()
		}
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			err := dec.conformantScan(v.Field(i), v.Type().Field(i).Tag)
			if err != nil {
				return err
			}
		}
	case reflect.String:
		ndrTag := parseTags(tag)
		if ndrTag.HasValue(TagPointer) || !ndrTag.HasValue(TagConformant) {
			break
		}
		dec.conformantMax = append(dec.conformantMax, uint32(0))
	case reflect.Slice:
		ndrTag := parseTags(tag)
		if ndrTag.HasValue(TagPointer) || !ndrTag.HasValue(TagConformant) {
			break
		}
		d, t := sliceDimensions(v.Type())
		for i := 0; i < d; i++ {
			dec.conformantMax = append(dec.conformantMax, uint32(0))
		}
		// For string arrays there is a common max for the strings within the array.
		if t.Kind() == reflect.String {
			dec.conformantMax = append(dec.conformantMax, uint32(0))
		}
	}
	return nil
}

func (dec *Decoder) fill(s interface{}, tag reflect.StructTag, deferred bool) error {
	var v reflect.Value
	if r, ok := s.(reflect.Value); ok {
		v = r
	} else {
		if reflect.ValueOf(s).Kind() == reflect.Ptr {
			v = reflect.ValueOf(s).Elem()
		}
	}
	ndrTag := parseTags(tag)
	// Pointer so defer filling the referent
	if ndrTag.HasValue(TagPointer) {
		dec.r.Discard(4) // discard the 4 bytes of the pointer
		ndrTag.delete(TagPointer)
		tag = ndrTag.StructTag()
		dec.deferred = append(dec.deferred, deferedPtr{v, tag})
		return nil
	}
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Implements(reflect.TypeOf(new(Union)).Elem()) {
			err := dec.readUnion(v, tag)
			if err != nil {
				return Errorf("could not fill union: %v", err)
			}
			break
		}
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Type().Implements(reflect.TypeOf(new(RawBytes)).Elem()) {
				//field is for rawbytes
				if v.Field(i).Type() != reflect.TypeOf([]byte{}) {
					return errors.New("cannot fill raw bytes as not a type of []byte")
				}
				err := dec.readRawBytes(v.Field(i))
				if err != nil {
					return Errorf("could not fill raw bytes: %v", err)
				}
			} else {
				err := dec.fill(v.Field(i), v.Type().Field(i).Tag, deferred)
				if err != nil {
					return Errorf("could not fill struct field(%d): %v", i, err)
				}
			}
		}
	case reflect.Bool:
		i, err := dec.readBool()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint8:
		i, err := dec.readUint8()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint16:
		i, err := dec.readUint16()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint32:
		i, err := dec.readUint32()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint64:
		i, err := dec.readUint64()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.String:
		ndrTag := parseTags(tag)
		conformant := ndrTag.HasValue(TagConformant)
		// strings are always varying so this is assumed
		var s string
		var err error
		if conformant {
			s, err = dec.readConformantVaryingString(deferred)
			if err != nil {
				return Errorf("could not fill with conformant varying string: %v", err)
			}
		} else {
			s, err = dec.readVaryingString()
			if err != nil {
				return Errorf("could not fill with varying string: %v", err)
			}
		}
		v.Set(reflect.ValueOf(s))
	case reflect.Float32:
		i, err := dec.readFloat32()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Float64:
		i, err := dec.readFloat64()
		if err != nil {
			return Errorf("could not fill %v: %v", v.Type().Name(), err)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Array:
		err := dec.fillFixedArray(v, tag)
		if err != nil {
			return err
		}
	case reflect.Slice:
		ndrTag := parseTags(tag)
		conformant := ndrTag.HasValue(TagConformant)
		varying := ndrTag.HasValue(TagVarying)
		_, t := sliceDimensions(v.Type())
		if t.Kind() == reflect.String && !ndrTag.HasValue(subStringArrayValue) {
			// String array
			err := dec.readStringsArray(v, tag, deferred)
			if err != nil {
				return err
			}
			break
		}
		// varying is assumed as fixed arrays use the Go array type rather than slice
		if conformant && varying {
			err := dec.fillConformantVaryingArray(v, tag, deferred)
			if err != nil {
				return err
			}
		} else if !conformant && varying {
			err := dec.fillVaryingArray(v, tag)
			if err != nil {
				return err
			}
		} else {
			//default to conformant and not varying
			err := dec.fillConformantArray(v, tag, deferred)
			if err != nil {
				return err
			}
		}
	default:
		return Errorf("unsupported type")
	}
	return nil
}

func (dec *Decoder) readBytes(n int) ([]byte, error) {
	//TODO make this take an int64 as input to allow for larger values on all systems?
	b := make([]byte, n, n)
	m, err := dec.r.Read(b)
	if err != nil || m != n {
		return b, fmt.Errorf("error reading bytes from stream: %v", err)
	}
	return b, nil
}
