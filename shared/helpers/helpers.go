package helpers

import "reflect"

func IsInterfaceStruct(intf interface{}) bool {
	return reflect.ValueOf(intf).Kind() != reflect.Struct
}

func HasStructField(structValue interface{}, fieldName string) bool {
	reflection := reflect.ValueOf(structValue)

	if reflection.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Type().Field(i)

		if field.Name == fieldName {
			return true
		}
	}

	return false
}
