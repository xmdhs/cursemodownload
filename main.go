package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/xmdhs/cursemodownload/web"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", web.Index)
	r.HandleFunc("/s", web.WebRoot)
	r.HandleFunc("/info", web.Info)
	s := http.Server{
		Addr:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	browser.OpenURL("http://127.0.0.1:8082")
	log.Println(s.ListenAndServe())
}
