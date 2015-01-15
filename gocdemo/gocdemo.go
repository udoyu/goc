package main

import (
	"flag"
	"fmt"
        "github.com/udoyu/goc/gochttp"
	"net/http"
        "github.com/xiying/xytool/simini"
	"time"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Printf("Usage:pp_server conf_file\n", flag.Arg(0))
		return
	}
	fmt.Println(flag.NArg(), flag.Arg(0))
	ini := simini.SimIni{}
	ini.LoadFile(flag.Arg(0))
	host := ini.GetStringVal("http-server", "host")
	fmt.Println("host:", host)
	time_out, _ := ini.GetIntValWithDefault("http-server", "reply_time_out", 3)
	gochttp.TimeoutDur = time.Second * time.Duration(time_out)
	strs := make([]string, flag.NArg()+1)
	strs[0] = "main"
	copy(strs[1:], flag.Args())
	fmt.Printf("%v\n", strs)
	gochttp.Init(len(strs), strs)
	HandlerStart(host)
}

func HandlerStart(port string) {
	myhandle := http.HandlerFunc(handlefun)
	err := http.ListenAndServe(port, myhandle)
	if nil != err {
		fmt.Println(err.Error())
	}
}

func handlefun(w http.ResponseWriter, req *http.Request) {
	gochttp.HttpEntry(req, &w)
}
