package helper

import (
	"fmt"
	"reflect"
)

// Reobject scannes the struct v
// and sets right attributes to s using json tag
func Reobject(src, dst interface{}, tagKey string) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst).Elem()

	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		fmt.Println(srcVal.Kind(), dstVal.Kind())
		return fmt.Errorf("both src and dst must be structs")
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		if tag, ok := srcField.Tag.Lookup(tagKey); ok {
			dstField := dstVal.FieldByName(tag)
			if dstField.IsValid() && dstField.CanSet() && dstField.Type() == srcVal.Field(i).Type() {
				dstField.Set(srcVal.Field(i))
			}
		}
	}

	return nil
}
