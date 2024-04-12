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
		return fmt.Errorf("both src and dst must be structs, got src: %s, dst: %s", srcVal.Kind(), dstVal.Kind())
	}
	srcFieldsByTag := make(map[string]reflect.Value)

	// Collect all source fields with their tag values
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		tag, ok := srcField.Tag.Lookup(tagKey)
		if ok {
			srcFieldsByTag[tag] = srcVal.Field(i)
		}
	}

	// Attempt to set destination fields based on matching tag values
	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Type().Field(i)
		tag, ok := dstField.Tag.Lookup(tagKey)
		if ok {
			if srcField, found := srcFieldsByTag[tag]; found {
				if dstVal.Field(i).CanSet() && srcField.Type() == dstVal.Field(i).Type() {
					dstVal.Field(i).Set(srcField)
				}
			}
		}
	}

	return nil
}
