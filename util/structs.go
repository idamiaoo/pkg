package util

import (
	"reflect"
)

func DeepFields(faceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < faceType.NumField(); i++ {
		v := faceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

func CopyStruct(copy interface{}, into interface{}) {
	copyValue := reflect.ValueOf(copy)
	intoValue := reflect.ValueOf(into)
	copyType := reflect.TypeOf(copy)
	intoType := reflect.TypeOf(into)
	if copyType.Kind() != reflect.Ptr || copyType.Elem().Kind() == reflect.Ptr || intoType.Kind() != reflect.Ptr || intoType.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}

	if copyValue.IsNil() || intoValue.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}

	copyV := copyValue.Elem()
	intoV := intoValue.Elem()
	fields := DeepFields(reflect.TypeOf(intoV.Interface()))

	for _, v := range fields {
		if v.Anonymous {
			continue
		}

		intoField := intoV.FieldByName(v.Name)
		copyField := copyV.FieldByName(v.Name)

		if !intoField.IsValid() {
			continue
		}

		if !copyField.IsValid() {
			continue
		}

		if copyField.Type() == intoField.Type() && intoField.CanSet() {
			intoField.Set(copyField)
			continue
		}

		if isKindInt(copyField.Type().Kind()) && isKindInt(intoField.Type().Kind()) && intoField.CanSet() {
			intoField.SetInt(copyField.Int())
			continue
		}

		if copyField.Kind() == reflect.Ptr && !copyField.IsNil() && copyField.Type().Elem() == copyField.Type() {
			intoField.Set(copyField.Elem())
			continue
		}

		if intoField.Kind() == reflect.Ptr && intoField.Type().Elem() == copyField.Type() {
			intoField.Set(reflect.New(copyField.Type()))
			intoField.Elem().Set(copyField)
			continue
		}
	}

	return
}

func isKindInt(kind reflect.Kind) bool {
	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int32 || kind == reflect.Int64 {
		return true
	}

	return false
}
