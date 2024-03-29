package dynamic

import (
	"fmt"
	"reflect"
)

// Call the field with the name "method" on the struct or map "on".
func Call(method string, on interface{}, idArg string) error {
	v := reflect.ValueOf(on)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		f := v.FieldByName(method)
		f.Call([]reflect.Value{reflect.ValueOf(idArg)})
	case reflect.Map:
		m := on.(map[string]interface{})
		m[method].(func(string))(idArg)
	default:
		return fmt.Errorf("unhandled Kind: %v(%v)", v.Kind(), v.Type())
	}
	return nil
}

// Ranger extracts a slice by its field name from an interface.
func Ranger(slice string, on interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(on)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	s := v.FieldByName(slice)
	if !s.IsValid() {
		return nil, fmt.Errorf("does not exist")
	}
	switch s.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		result := make([]interface{}, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			result = append(result, s.Index(i).Interface())
		}
		return result, nil
	default:
		return nil, fmt.Errorf("kind not iterable: %v", s.Kind())
	}
}
