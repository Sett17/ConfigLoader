package fieldsetter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SetFields updates the fields of the given object based on a map of field paths to values.
// Field names in the path must exactly match the struct field names, including capitalization.
// It optionally continues on error when soft is true, collecting all errors encountered.
// Returns a slice of errors encountered during the update, or nil if no errors occurred.
func SetFields(obj any, fields map[string]any, soft bool) []error {
	var errs []error
	for path, value := range fields {
		err := SetValue(obj, path, value)
		if err != nil {
			if soft {
				errs = append(errs, err)
			} else {
				return []error{err}
			}
		}
	}
	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}

// SetValue updates a specific field of an object based on the given field path and value.
// The object must be a pointer, and the field path must exactly match the struct's field names,
// including capitalization. Supports nested fields, arrays, slices, and maps. For arrays and slices,
// the path segment should include the index directly following the field name (e.g., "ArrayField.0").
// For maps, it should include the key directly following the field name (e.g., "MapField.Key").
// Returns an error if the object is not a pointer, the path is invalid, the specified index is out of
// bounds, the key does not exist in the map, or if the value type is incompatible with the field,
// array element, or map value type.
func SetValue(obj any, path string, value any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("Object must be a pointer")
	}
	pathSegments := strings.Split(path, ".")
	return setFieldRecursive(v, pathSegments, value)
}

func setFieldRecursive(v reflect.Value, pathSegments []string, value any) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	if v.Kind() != reflect.Pointer || !v.Elem().CanSet() {
		return errors.New("target must be a pointer and settable")
	}
	v = v.Elem()

	switch v.Kind() {
	case reflect.Struct:
		field := v.FieldByName(pathSegments[0])
		if !field.IsValid() {
			return fmt.Errorf("field %s does not exist", pathSegments[0])
		}
		if !field.CanSet() {
			return fmt.Errorf("cannot set field %s", pathSegments[0])
		}
		if len(pathSegments) == 1 {
			if value == nil {
				field.Set(reflect.Zero(field.Type()))
				return nil
			}
			fieldValue := reflect.ValueOf(value)
			if fieldValue.Type().AssignableTo(field.Type()) {
				field.Set(fieldValue)
				return nil
			}
			return fmt.Errorf("value type %s is not assignable to field type %s", fieldValue.Type(), field.Type())
		}
		return setFieldRecursive(field.Addr(), pathSegments[1:], value)
	case reflect.Slice, reflect.Array:
		if len(pathSegments) < 1 {
			return fmt.Errorf("no index provided for slice/array")
		}
		index, err := strconv.Atoi(pathSegments[0])
		if err != nil {
			return fmt.Errorf("invalid index: %s", pathSegments[0])
		}
		if index < 0 || index >= v.Len() {
			return fmt.Errorf("index out of range: %d", index)
		}
		if len(pathSegments) == 1 {
			elem := v.Index(index)
			newValue := reflect.ValueOf(value)
			if newValue.Type().AssignableTo(elem.Type()) {
				elem.Set(newValue)
				return nil
			}
			return fmt.Errorf("value type %s is not assignable to element type %s", newValue.Type(), elem.Type())
		}
		return setFieldRecursive(v.Index(index).Addr(), pathSegments[1:], value)
	case reflect.Map:
		if len(pathSegments) < 1 {
			return fmt.Errorf("no key provided for map")
		}
		key := reflect.ValueOf(pathSegments[0])
		newValue := reflect.ValueOf(value)
		if !newValue.Type().AssignableTo(v.Type().Elem()) {
			return fmt.Errorf("value type %s is not assignable to map value type %s", newValue.Type(), v.Type().Elem())
		}
		v.SetMapIndex(key, newValue)
		return nil
	default:
		return fmt.Errorf("unsupported type %s", v.Kind())
	}
}
