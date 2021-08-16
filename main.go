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
	r.HandleFunc("/curseforge", web.Index)
	r.HandleFunc("/curseforge/s", web.WebRoot)
	r.HandleFunc("/curseforge/info", web.Info)
	r.HandleFunc("/curseforge/download", web.Getdownloadlink)
	r.HandleFunc("/curseforge/history", web.History)
	s := http.Server{
		Addr:         "127.23.51.3:80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	browser.OpenURL("http://127.23.51.3:80/curseforge")
	log.Println(s.ListenAndServe())
}
