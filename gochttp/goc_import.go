package gochttp

import "C"
import (
	"unsafe"
)

//export SendData
func SendData(ctx unsafe.Pointer, data *C.char, length C.int) {
	c := (*chan_data_t)(ctx)
	if c.is_disable {
		return
	}
	var strptr *string = nil
	if length > 0 {
		str := C.GoStringN(data, length)
		strptr = &str
	}
	c.Add(&str_t{strptr, 0}, 0)
}

//export Finalize
func Finalize(ctx unsafe.Pointer, result C.int) {
	c := (*chan_data_t)(ctx)
	if c.is_disable {
		return
	}
	c.Add(&str_t{nil, int(result)}, 0)
}
