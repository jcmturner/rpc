package ndr

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	defaultAlignment = 4
	TagConformant    = "conformant"
	TagVarying       = "varying"
)

type Decoder struct {
	//mutex sync.Mutex    // each item must be received atomically
	r     *bufio.Reader // source of the data
	size  int           // initial size of bytes in buffer
	align int           // the alignment multiple
	ch    CommonHeader  // NDR common header
	ph    PrivateHeader // NDR private header
}

func NewDecoder(r io.Reader, align int) *Decoder {
	dec := new(Decoder)
	dec.r = bufio.NewReader(r)
	dec.size = dec.r.Buffered()
	if align != 1 && align != 2 && align != 4 && align != 8 {
		align = defaultAlignment
	}
	dec.align = align
	return dec
}

func (dec *Decoder) Decode(s interface{}) error {
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
	return dec.fill(s, reflect.StructTag(""))
}

func (dec *Decoder) fill(s interface{}, tag reflect.StructTag) error {
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
			dec.fill(v.Field(i), v.Type().Field(i).Tag)
		}
	case reflect.Bool:
		i, err := dec.readBool()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint8:
		i, err := dec.readUint8()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint16:
		i, err := dec.readUint16()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint32:
		i, err := dec.readUint32()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Uint64:
		i, err := dec.readUint64()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.String:
		ndrTag := parseTags(tag)
		conformant := ndrTag.HasValue(TagConformant)
		// strings are always varying so this is assumed
		var s string
		var err error
		if conformant {
			s, err = dec.readConformantVaryingString()
			if err != nil {
				return fmt.Errorf("could not fill with conformant varying string %v", v)
			}
		} else {
			s, err = dec.readVaryingString()
			if err != nil {
				return fmt.Errorf("could not fill with varying string %v", v)
			}
		}
		v.Set(reflect.ValueOf(s))
	case reflect.Slice:

	case reflect.Float32:
		i, err := dec.readFloat32()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Float64:
		i, err := dec.readFloat64()
		if err != nil {
			return fmt.Errorf("could not fill %v", v)
		}
		v.Set(reflect.ValueOf(i))
	default:
		return errors.New("unsupported type")
	}
	return nil
}

func (dec *Decoder) ensureAlignment() {
	p := dec.size - dec.r.Buffered()
	if s := p % dec.align; s != 0 {
		dec.r.Discard(dec.align - s)
	}
}

func (dec *Decoder) readBytes(n int) ([]byte, error) {
	b := make([]byte, n, n)
	m, err := dec.r.Read(b)
	if err != nil || m != n {
		return b, Malformed{EText: fmt.Sprintf("could not read bytes from stream: %v", err)}
	}
	return b, nil
}
