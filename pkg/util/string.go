package util

import (
	"regexp"
	"strconv"
	"unsafe"
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func PhoneVerify(phone string) bool {
	reg := `^1[3456789][0-9]{9}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func ByteUintToint64(bs []uint8) int64 {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	in, _ := strconv.Atoi(string(ba))
	return int64(in)
}

func ByteUintToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
