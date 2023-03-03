package renderer

import (
	"reflect"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

func resolveCollisions(value interface{}) {
	resolve(reflect.ValueOf(value))
}

func changeSlice(rv reflect.Value) {
	if !rv.CanAddr() {
		return
	}

	for i := 0; i < rv.Len(); i++ {
		forEachSliceElem(rv.Index(i))
	}
}

func resolve(rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		changeStruct(rv)
	case reflect.Slice:
		changeSlice(rv)
	}
}

func changeStruct(rv reflect.Value) {
	if !rv.CanAddr() {
		return
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		forEachField(field)
	}
}

func forEachField(field reflect.Value) {
	switch field.Kind() {
	case reflect.String:
		newValue := strutils.ResolveNameIfCollides(field.String())
		field.SetString(newValue)
	case reflect.Slice:
		changeSlice(field)
	}
}

func forEachSliceElem(item reflect.Value) {
	switch item.Kind() {
	case reflect.Struct:
		resolve(item)
	case reflect.String:
		newValue := strutils.ResolveNameIfCollides(item.String())
		item.SetString(newValue)
	}
}
