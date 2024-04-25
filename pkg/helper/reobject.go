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

	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Type().Field(i)
		tag, ok := dstField.Tag.Lookup(tagKey)
		if ok {
			if srcField, found := srcFieldsByTag[tag]; found {
				dstFieldVal := dstVal.Field(i)
				if dstFieldVal.CanSet() {
					if srcField.Kind() == reflect.Ptr && dstFieldVal.Kind() != reflect.Ptr {
						if !srcField.IsNil() && srcField.Elem().Type() == dstFieldVal.Type() {
							// Dereference the pointer if src is not nil and types match
							dstFieldVal.Set(srcField.Elem())
						}
					} else if srcField.Type() == dstFieldVal.Type() {
						// Set directly if the types are exactly the same
						dstFieldVal.Set(srcField)
					}
				}
			}
		}
	}

	return nil
}
