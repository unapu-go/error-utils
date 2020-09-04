package error_utils

import "reflect"

func WalkErr(cb func(err error) (stop bool), errs ...error) (stop bool) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		if cb(err) {
			return true
		}

		if err, ok := err.(interface{ Err() error }); ok {
			if WalkErr(cb, err.Err()) {
				return true
			}
		}
		if err, ok := err.(Causer); ok {
			if WalkErr(cb, err.Cause()) {
				return true
			}
		}

		if errs, ok := err.(Errors); ok {
			if WalkErr(cb, errs...) {
				return true
			}
		} else if errs, ok := err.(interface{ Errors() []error }); ok {
			if WalkErr(cb, errs.Errors()...) {
				return true
			}
		} else if errs, ok := err.(interface{ GetErrors() []error }); ok {
			if WalkErr(cb, errs.GetErrors()...) {
				return true
			}
		}
	}
	return false
}

func IsError(expected error, err ...error) (is bool) {
	return WalkErr(func(err error) (stop bool) {
		return err == expected
	}, err...)
}

func ErrorByType(expected reflect.Type, err ...error) (theError error) {
	expected = indirectRealType(expected)
	WalkErr(func(err error) (stop bool) {
		if indirectRealType(reflect.TypeOf(err)) == expected {
			theError = err
			return true
		}
		return false
	}, err...)
	return
}

func ErrorByInterfaceType(expected reflect.Type, err ...error) (theError error) {
	for expected.Kind() == reflect.Ptr {
		expected = expected.Elem()
	}
	WalkErr(func(err error) (stop bool) {
		if reflect.TypeOf(err).Implements(expected) {
			theError = err
			return true
		}
		return false
	}, err...)
	return
}

func IsErrorTyp(expected reflect.Type, err ...error) (is bool) {
	return ErrorByType(expected, err...) != nil
}

func TraceOf(err ...error) (st []byte) {
	if terr := ErrorByInterfaceType(reflect.TypeOf((*Tracer)(nil)), err...); terr != nil {
		st = terr.(Tracer).Trace()
	}
	return
}