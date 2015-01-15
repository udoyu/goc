package gochttp

import (
	"C"
        "github.com/astaxie/beego/session"
	"unsafe"
)

//export SetSessInt
func SetSessInt(psess unsafe.Pointer, name *C.char, value C.int) {
	sess := (*session.SessionStore)(psess)
	(*sess).Set(C.GoString(name), int(value))
}

//export SetSessStr
func SetSessStr(psess unsafe.Pointer, name *C.char, value *C.char) {
	sess := (*session.SessionStore)(psess)
	(*sess).Set(C.GoString(name), C.GoString(value))
}

//export GetSessInt
func GetSessInt(psess unsafe.Pointer, name *C.char) C.int {
	sess := (*session.SessionStore)(psess)
	v := (*sess).Get(C.GoString(name))
	if v == nil {
		return 0
	}
	return C.int(v.(int))
}

//export GetSessStr
func GetSessStr(psess unsafe.Pointer, name *C.char) *C.char {
	sess := (*session.SessionStore)(psess)
	v := (*sess).Get(C.GoString(name))
	if v == nil {
		return C.CString("")
	}
	return C.CString(v.(string))
}

//export DelSess
func DelSess(psess unsafe.Pointer, name *C.char) {
	sess := (*session.SessionStore)(psess)
	(*sess).Delete(C.GoString(name))
}

//export FlushSess
func FlushSess(psess unsafe.Pointer) {
	sess := (*session.SessionStore)(psess)
	(*sess).Flush()
}
