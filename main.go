package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "unsafe"

	"github.com/pkg/browser"
	_ "github.com/xmdhs/cursemodownload/curseapi"
	"github.com/xmdhs/cursemodownload/web"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/curseforge", web.Index)
	r.HandleFunc("/curseforge/s", web.WebRoot)
	r.HandleFunc("/curseforge/info", web.Info)
	r.HandleFunc("/curseforge/download", web.Getdownloadlink)
	r.HandleFunc("/curseforge/history", web.History)
	s := http.Server{
		Addr:         "127.0.0.1:8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	fmt.Println("WebServer Starting...")
	browser.OpenURL("http://127.0.0.1:8082/curseforge")
	log.Println(s.ListenAndServe())
}

var apiaddr string

//go:linkname api github.com/xmdhs/cursemodownload/curseapi.api
var api string

func init() {
	flag.StringVar(&apiaddr, "apiaddr", "https://addons-ecs.forgesvc.net", "api address")
	flag.Parse()

	api = apiaddr
}
