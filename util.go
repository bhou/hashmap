package hashmap

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/dchest/siphash"
)

// intSizeBytes is the size in byte of an int or uint value.
const intSizeBytes = strconv.IntSize >> 3

// roundUpPower2 rounds a number to the next power of 2.
func roundUpPower2(i uint64) uint64 {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}

// log2 computes the binary logarithm of x, rounded up to the next integer.
func log2(i uint64) uint64 {
	var n, p uint64
	for p = 1; p < i; p += p {
		n++
	}
	return n
}

func getKeyHash(key interface{}) uint64 {
	var num uint64
	v := reflect.ValueOf(key)

	switch reflect.TypeOf(key).Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
		return 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = uint64(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		num = uint64(v.Uint())

	case reflect.String:
		s := key.(string)
		sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
		bh := reflect.SliceHeader{
			Data: sh.Data,
			Len:  sh.Len,
			Cap:  sh.Len,
		}
		buf := *(*[]byte)(unsafe.Pointer(&bh))
		return siphash.Hash(1, 2, buf)

	default:
		panic("unsupported key type")
	}

	bh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&num)),
		Len:  8,
		Cap:  8,
	}
	buf := *(*[]byte)(unsafe.Pointer(&bh))
	return siphash.Hash(1, 2, buf)
}
