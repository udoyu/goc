package gochttp

//#cgo LDFLAGS: -lhttpentry
//#include "http_entry.h"
/*
#include <stdio.h>
extern void SendData(void* ctx, char* data, int length);
extern void Finalize(void* ctx, int result);
extern void SetSessInt(void* sess, char* name, int value);
extern void SetSessStr(void* sess, char* name, char* value);
extern int GetSessInt(void* sess, char* name);
extern char* GetSessStr(void* sess, char* name);
extern void DelSess(void* sess, char* name);
extern void FlushSess(void* sess);
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"net/http"
        "github.com/astaxie/beego/session"
	"time"
	"unsafe"
)

var TimeoutDur time.Duration
var init_flag bool
var globalSessions *session.Manager

func init() {
	TimeoutDur = time.Second * time.Duration(3)
	//globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": -1, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 0, "providerConfig": "127.0.0.1:6379"}`)
	//if globalSessions == nil {
	//	fmt.Println("sess is nil")
	//	return
	//}
	//go globalSessions.GC()
}

func InitSession(s *session.Manager) {
    globalSessions = s
}

func Init(argc int, argv []string) {
	if !init_flag {
		fmt.Println("argc:", argc)
		argvs := make([]*C.char, argc)
		for i, v := range argv {
			argvs[i] = C.CString(v)
		}
		C.Init(C.int(argc), &argvs[0])
		init_flag = true
	}
}

func Exit(argc int, argv []string) {
	argvs := make([]*C.char, argc)
	for i, v := range argv {
		argvs[i] = C.CString(v)
	}
	C.Exit(C.int(argc), &argvs[0])
}

func HttpEntry(req *http.Request, resp *http.ResponseWriter) {
	var httpreq C.http_request_t
	httpreq.uri = C.CString(req.URL.Path)
	switch req.Method {
	case "GET":
		httpreq.method = C.GET
	case "POST":
		httpreq.method = C.POST
	default:
		httpreq.method = C.OTHER
	}
	req.ParseForm()
	//	req.FormValue("key")
	httpreq.arglen = C.size_t(len(req.Form))
	args := make([]C.request_arg_t, len(req.Form))
	i := 0
	for key, value := range req.Form {
		args[i].name = C.CString(key)
		args[i].value.length = C.size_t(len(value[0]))
		args[i].value.data = C.CString(value[0])
		i++
	}
	if len(args) > 0 {
		httpreq.args = &args[0]
	}
	body, _ := ioutil.ReadAll(req.Body)
	if len(body) > 0 {
		bodystr := string(body)
		httpreq.body.length = C.size_t(len(bodystr))
		httpreq.body.data = C.CString(bodystr)
	}
	//创建buf,buf[0]存储要发送的数据,buf[1]存储c方法返回的数据
	data := NewChanStr(1)
	httpreq.ctx = unsafe.Pointer(data)
	httpreq.sendfunc = C.send_data_pt(C.SendData)
	httpreq.finalizefunc = C.finalize_request_pt(C.Finalize)
	sess := globalSessions.SessionStart(*resp, req)
	defer sess.SessionRelease(*resp)
	httpreq.sess.sess = unsafe.Pointer(&sess)
	httpreq.sess.set_int = C.set_sess_int_pt(C.SetSessInt)
	httpreq.sess.set_str = C.set_sess_str_pt(C.SetSessStr)
	httpreq.sess.get_int = C.get_sess_int_pt(C.GetSessInt)
	httpreq.sess.get_str = C.get_sess_str_pt(C.GetSessStr)
	httpreq.sess.del = C.del_sess_pt(C.DelSess)
	httpreq.sess.flush = C.flush_sess_pt(C.FlushSess)
	go C.HttpEntry(&httpreq)
	timeout := time.NewTicker(TimeoutDur)
	str := ""
	strptr := &str
	select {
	case senddata := <-data.data[0]:
		strptr = senddata.data
		break
	case <-timeout.C:
		data.is_disable = true
		*strptr = "timeout"
		break
	}
	if nil != strptr {
		fmt.Fprintf(*resp, *strptr)
	}
}
