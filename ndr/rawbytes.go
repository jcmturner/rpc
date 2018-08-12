package ndr

import (
	"errors"
	"fmt"
	"reflect"
)

// type MyBytes []byte
// implement RawBytes interface

const (
	sizeMethod = "Size"
)

type RawBytes interface {
	Size(interface{}) int
}

func (dec *Decoder) readRawBytes(v reflect.Value) error {
	sf := v.MethodByName(sizeMethod)
	if !sf.IsValid() {
		return fmt.Errorf("could not find a method called %s on the implementation of RawBytes", sizeMethod)
	}
	in := []reflect.Value{reflect.ValueOf(dec.s)}
	f := sf.Call(in)
	if f[0].Kind() != reflect.Int {
		return errors.New("the RawBytes size function did not return an integer")
	}
	b, err := dec.readBytes(int(f[0].Int()))
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(b).Convert(v.Type()))
	return nil
}
