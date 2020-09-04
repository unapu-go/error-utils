package error_utils

import "reflect"

func indirectRealType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Interface {
		reflectType = reflectType.Elem()
	}
	return reflectType
}
