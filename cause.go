package error_utils

type Causer interface {
	error
	Cause() error
}
