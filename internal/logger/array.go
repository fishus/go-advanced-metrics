package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Array constructs a field with the given key and ArrayMarshaler. It provides
// a flexible, but still type-safe and efficient, way to add array-like types
// to the logging context. The struct's MarshalLogArray method is called lazily.
func Array(key string, val zapcore.ArrayMarshaler) Field {
	return zap.Array(key, val)
}

// Bools constructs a field that carries a slice of bools.
func Bools(key string, bs []bool) Field {
	return zap.Bools(key, bs)
}

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, bss [][]byte) Field {
	return zap.ByteStrings(key, bss)
}

// Complex128s constructs a field that carries a slice of complex numbers.
func Complex128s(key string, nums []complex128) Field {
	return zap.Complex128s(key, nums)
}

// Complex64s constructs a field that carries a slice of complex numbers.
func Complex64s(key string, nums []complex64) Field {
	return zap.Complex64s(key, nums)
}

// Durations constructs a field that carries a slice of time.Durations.
func Durations(key string, ds []time.Duration) Field {
	return zap.Durations(key, ds)
}

// Float64s constructs a field that carries a slice of floats.
func Float64s(key string, nums []float64) Field {
	return zap.Float64s(key, nums)
}

// Float32s constructs a field that carries a slice of floats.
func Float32s(key string, nums []float32) Field {
	return zap.Float32s(key, nums)
}

// Ints constructs a field that carries a slice of integers.
func Ints(key string, nums []int) Field {
	return zap.Ints(key, nums)
}

// Int64s constructs a field that carries a slice of integers.
func Int64s(key string, nums []int64) Field {
	return zap.Int64s(key, nums)
}

// Int32s constructs a field that carries a slice of integers.
func Int32s(key string, nums []int32) Field {
	return zap.Int32s(key, nums)
}

// Int16s constructs a field that carries a slice of integers.
func Int16s(key string, nums []int16) Field {
	return zap.Int16s(key, nums)
}

// Int8s constructs a field that carries a slice of integers.
func Int8s(key string, nums []int8) Field {
	return zap.Int8s(key, nums)
}

// Objects constructs a field with the given key, holding a list of the
// provided objects that can be marshaled by Zap.
//
// Note that these objects must implement zapcore.ObjectMarshaler directly.
// That is, if you're trying to marshal a []Request, the MarshalLogObject
// method must be declared on the Request type, not its pointer (*Request).
// If it's on the pointer, use ObjectValues.
//
// Given an object that implements MarshalLogObject on the value receiver, you
// can log a slice of those objects with Objects like so:
//
//	type Author struct{ ... }
//	func (a Author) MarshalLogObject(enc zapcore.ObjectEncoder) error
//
//	var authors []Author = ...
//	logger.Info("loading article", zap.Objects("authors", authors))
//
// Similarly, given a type that implements MarshalLogObject on its pointer
// receiver, you can log a slice of pointers to that object with Objects like
// so:
//
//	type Request struct{ ... }
//	func (r *Request) MarshalLogObject(enc zapcore.ObjectEncoder) error
//
//	var requests []*Request = ...
//	logger.Info("sending requests", zap.Objects("requests", requests))
//
// If instead, you have a slice of values of such an object, use the
// ObjectValues constructor.
//
//	var requests []Request = ...
//	logger.Info("sending requests", zap.ObjectValues("requests", requests))
func Objects[T zapcore.ObjectMarshaler](key string, values []T) Field {
	return zap.Objects(key, values)
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// Stringers constructs a field with the given key, holding a list of the
// output provided by the value's String method
//
// Given an object that implements String on the value receiver, you
// can log a slice of those objects with Objects like so:
//
//	type Request struct{ ... }
//	func (a Request) String() string
//
//	var requests []Request = ...
//	logger.Info("sending requests", zap.Stringers("requests", requests))
//
// Note that these objects must implement fmt.Stringer directly.
// That is, if you're trying to marshal a []Request, the String method
// must be declared on the Request type, not its pointer (*Request).
func Stringers[T fmt.Stringer](key string, values []T) Field {
	return zap.Stringers(key, values)
}

// Times constructs a field that carries a slice of time.Times.
func Times(key string, ts []time.Time) Field {
	return zap.Times(key, ts)
}

// Uints constructs a field that carries a slice of unsigned integers.
func Uints(key string, nums []uint) Field {
	return zap.Uints(key, nums)
}

// Uint64s constructs a field that carries a slice of unsigned integers.
func Uint64s(key string, nums []uint64) Field {
	return zap.Uint64s(key, nums)
}

// Uint32s constructs a field that carries a slice of unsigned integers.
func Uint32s(key string, nums []uint32) Field {
	return zap.Uint32s(key, nums)
}

// Uint16s constructs a field that carries a slice of unsigned integers.
func Uint16s(key string, nums []uint16) Field {
	return zap.Uint16s(key, nums)
}

// Uint8s constructs a field that carries a slice of unsigned integers.
func Uint8s(key string, nums []uint8) Field {
	return zap.Uint8s(key, nums)
}

// Uintptrs constructs a field that carries a slice of pointer addresses.
func Uintptrs(key string, us []uintptr) Field {
	return zap.Uintptrs(key, us)
}

// Errors constructs a field that carries a slice of errors.
func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}
