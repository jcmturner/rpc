package ndr

import "reflect"

func parseDimensions(v reflect.Value, t reflect.Type) (d []reflect.Kind) {
	if t.Elem().Kind() == reflect.Array {
		d = append(d, parseDimensions(v, t)...) //TODO new values of v and t
	}
	return
}
