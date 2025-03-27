package helper

import (
	"unsafe"
)

func BytesInt[T ~int](i T) []byte {
	return (*[8]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt8[T ~int8](i T) []byte {
	return (*[1]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt16[T ~int16](i T) []byte {
	return (*[2]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt32[T ~int32](i T) []byte {
	return (*[4]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt64[T ~int64](i T) []byte {
	return (*[8]byte)(unsafe.Pointer(&i))[:]
}

func BytesFloat32[T ~float32](f T) []byte {
	return (*[4]byte)(unsafe.Pointer(&f))[:]
}

func BytesFloat64[T ~float64](f T) []byte {
	return (*[8]byte)(unsafe.Pointer(&f))[:]
}

func BytesString[T ~string](s T) []byte {
	return unsafe.Slice(unsafe.StringData(string(s)), len(s))
}

func BytesBool[T ~bool](b T) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}
