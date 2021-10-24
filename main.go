package main

import (
	"log"
	"net/http"
	"time"
	"fmt"

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
		Addr:         "127.0.0.1:11451",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	fmt.Println("WebServer Starting...")
	browser.OpenURL("http://localhost:11451/curseforge")
	log.Println(s.ListenAndServe())
}
