package helper

import (
	"unsafe"
)

func BytesInt(i int) []byte {
	return (*[8]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt8(i int8) []byte {
	return (*[1]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt16(i int16) []byte {
	return (*[2]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt32(i int32) []byte {
	return (*[4]byte)(unsafe.Pointer(&i))[:]
}

func BytesInt64(i int64) []byte {
	return (*[8]byte)(unsafe.Pointer(&i))[:]
}

func BytesFloat32(f float32) []byte {
	return (*[4]byte)(unsafe.Pointer(&f))[:]
}

func BytesFloat64(f float64) []byte {
	return (*[8]byte)(unsafe.Pointer(&f))[:]
}

func BytesString(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}
