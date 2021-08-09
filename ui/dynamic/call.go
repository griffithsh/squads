package dynamic

import (
	"fmt"
	"reflect"
)

// Call the field with the name "method" on the struct or map "on".
func Call(method string, on interface{}) error {
	v := reflect.ValueOf(on)
	switch v.Kind() {
	case reflect.Struct:
		f := v.FieldByName(method)
		f.Call([]reflect.Value{})
	case reflect.Map:
		m := on.(map[string]interface{})
		m[method].(func())()
	default:
		return fmt.Errorf("unhandled Kind: %v(%v)", v.Kind(), v.Type())
	}
	return nil
}
