package ndr

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
)

const (
	defaultAlignment = 4
)

type Decoder struct {
	mutex  sync.Mutex    // each item must be received atomically
	r      *bufio.Reader // source of the data
	size   int           //initial size of bytes in buffer
	align  int
	ch     CommonHeader
	ph     PrivateHeader
	fields []field
	err    error
}

type field struct {
	Name string
	Type string
	Tags string
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
	dec.fill(s)
	return nil
}

func (dec *Decoder) fill(s interface{}) error {
	v := reflect.ValueOf(&s).Elem()
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			dec.fill(v.Field(i))
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

	case reflect.Slice:

	case reflect.Float32:

	case reflect.Float64:

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
