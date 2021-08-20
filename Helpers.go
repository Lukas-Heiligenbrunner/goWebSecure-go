package goWebSecure_go

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Jsonify(v interface{}) []byte {
	// jsonify results
	str, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error while Jsonifying return object: " + err.Error())
	}
	return str
}

// setField set a specific field of an object with an object provided
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType != val.Type() {
		if val.Type().ConvertibleTo(structFieldType) {
			// if type is convertible - convert and set
			structFieldValue.Set(val.Convert(structFieldType))
		} else {
			return fmt.Errorf("provided value %s type didn't match obj field type and isn't convertible", name)
		}
	} else {
		// set value if type is the same
		structFieldValue.Set(val)
	}

	return nil
}

// FillStruct fill a custom struct with objects of a map
func FillStruct(i interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(i, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
