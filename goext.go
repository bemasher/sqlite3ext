package main

/*
#cgo CFLAGS: -std=gnu99
#cgo CFLAGS: -DUSE_LIBSQLITE3
#cgo LDFLAGS: -lsqlite3

#include "sqlite3.h"
#include <stdlib.h>
*/
import "C"

import (
	"log"
	"regexp"
	"time"
	"unsafe"

	"github.com/bemasher/stringdist"
)

const valueSliceLength = (1<<31 - 1) / unsafe.Sizeof((*C.struct_sqlite3_value)(nil))

func sqlite3_value_blob(v *C.sqlite3_value) []byte {
	n := C.sqlite3_value_bytes(v)
	b := C.sqlite3_value_blob(v)
	return C.GoBytes(b, n)
}

func sqlite3_value_text(v *C.sqlite3_value) string {
	n := C.sqlite3_value_bytes(v)
	b := (*C.char)(unsafe.Pointer(C.sqlite3_value_text(v)))
	return C.GoStringN(b, n)
}

func sqlite3_value_int(v *C.sqlite3_value) int {
	return int(C.sqlite3_value_int(v))
}

func sqlite3_value_bool(v *C.sqlite3_value) bool {
	return int(C.sqlite3_value_int(v)) != 0
}

func sqlite3_result_error(ctx *C.struct_sqlite3_context, err error) {
	cerr := C.CString(err.Error())
	C.sqlite3_result_error(ctx, cerr, -1)
	C.free(unsafe.Pointer(cerr))
}

func sqlite3_result_text(ctx *C.struct_sqlite3_context, text string) {
	C.sqlite3_result_text(ctx, C.CString(text), -1, (*[0]byte)(unsafe.Pointer(C.free)))
}

//export Jaro
func Jaro(ctx *C.struct_sqlite3_context, argc C.int, argv **C.struct_sqlite3_value) {
	args := (*[valueSliceLength]*C.sqlite3_value)(unsafe.Pointer(argv))[:argc:argc]

	s1 := sqlite3_value_text(args[0])
	s2 := sqlite3_value_text(args[1])

	C.sqlite3_result_double(ctx, C.double(stringdist.Jaro(s1, s2)))
}

var reCache map[string]*regexp.Regexp

//export Regex
func Regex(ctx *C.struct_sqlite3_context, argc C.int, argv **C.struct_sqlite3_value) {
	args := (*[valueSliceLength]*C.sqlite3_value)(unsafe.Pointer(argv))[:argc:argc]

	pattern := sqlite3_value_text(args[0])
	value := sqlite3_value_text(args[1])

	if _, ok := reCache[pattern]; !ok {
		reCache[pattern] = regexp.MustCompile(pattern)
	}

	match := reCache[pattern].MatchString(value)
	if match {
		C.sqlite3_result_int(ctx, C.int(1))
	} else {
		C.sqlite3_result_int(ctx, C.int(0))
	}
}

//export ParseTime
func ParseTime(ctx *C.struct_sqlite3_context, argc C.int, argv **C.struct_sqlite3_value) {
	args := (*[valueSliceLength]*C.sqlite3_value)(unsafe.Pointer(argv))[:argc:argc]

	srcFmt := sqlite3_value_text(args[0])
	dstFmt := sqlite3_value_text(args[1])
	src := sqlite3_value_text(args[2])
	grace := sqlite3_value_bool(args[3])
	dst, err := parseTime(srcFmt, dstFmt, src)
	if err != nil {
		if grace {
			sqlite3_result_text(ctx, src)
		} else {
			sqlite3_result_error(ctx, err)
		}
		return
	}

	sqlite3_result_text(ctx, dst)
}

func parseTime(srcFmt, dstFmt, val string) (dst string, err error) {
	t, err := time.Parse(srcFmt, val)
	if err != nil {
		return "", err
	}

	switch dstFmt {
	case "date":
		dstFmt = "2006-01-02"
	case "time":
		dstFmt = "15:04:05"
	case "datetime":
		dstFmt = "2006-01-02 15:04:05"
	}

	return t.Format(dstFmt), nil
}

func init() {
	reCache = make(map[string]*regexp.Regexp)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {}
