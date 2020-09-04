package error_utils

type Tracer interface {
	error
	Trace() []byte
}
