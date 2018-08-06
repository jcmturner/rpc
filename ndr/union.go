package ndr

import (
	"errors"
	"fmt"
	"reflect"
)

// To read a union a struct must be defined that has a field called "Tag" that is of the correct type (integer, char or
// bool) for the union's descriminating tag value.
// The struct must also have a method called SwitchFunc that takes the discriminating tag as input and returns a string
// for the name of the field that should be filled with the selected value.
// Below is a basic example:
//
// type MyUnion struct {
// 	 Tag uint32 // Can be any type but must be called Tag
// 	 Value32 uint32
//	 Value64 uint64
// }
//
// // Note that this method does not have a pointer receiver.
// func (u MyUnion) SwitchFunc(tag uint32) string {
//	 switch tag {
//	 case 64:
//		 return "Value64"
//	 case 32:
//		 return "Value32"
//	 }
//	 return ""
// }

const (
	DiscriminatingTagFieldName = "Tag"
	SelectionFuncName          = "SwitchFunc"
	TagEncapsulated            = "encapsulated"
	TagUnion                   = "union"
)

func (dec *Decoder) readUnion(v reflect.Value, tag reflect.StructTag) error {
	var utTag reflect.StructTag
	if utt, ok := v.Type().FieldByName(DiscriminatingTagFieldName); ok {
		utTag = utt.Tag
	} else {
		return errors.New("could not get union's discriminating Tag field")
	}
	ut := v.FieldByName(DiscriminatingTagFieldName)
	ndrTag := parseTags(tag)
	// For a non-encapsulated union, the discriminant is marshalled into the transmitted data stream twice: once as the
	// field or parameter, which is referenced by the switch_is construct, in the procedure argument list; and once as
	// the first part of the union representation.
	if !ndrTag.HasValue(TagEncapsulated) {
		dec.r.Discard(int(ut.Type().Size()))
	}
	err := dec.fill(ut, utTag)
	if err != nil {
		return fmt.Errorf("could not fill union's discriminating tag: %v", err)
	}
	args := []reflect.Value{ut}
	// Call the SelectFunc of the union struct to find the name of the field to fill with the value selected.
	sf := v.MethodByName(SelectionFuncName)
	if !sf.IsValid() {
		return fmt.Errorf("could not find a selection function called %s in the unions struct representation", SelectionFuncName)
	}
	f := sf.Call(args)
	if f[0].Kind() != reflect.String && f[0].String() != "" {
		return fmt.Errorf("the union select function did not return a string for the name of the field to fill")
	}

	var uvTag reflect.StructTag
	if uvt, ok := v.Type().FieldByName(f[0].String()); ok {
		uvTag = uvt.Tag
	} else {
		return fmt.Errorf("could not get union's selected value field: %s", f[0].String())
	}
	uv := v.FieldByName(f[0].String())
	err = dec.fill(uv, uvTag)
	if err != nil {
		return err
	}
	return nil
}
