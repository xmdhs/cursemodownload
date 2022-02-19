package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type historyS struct {
	Name             string
	Version          string
	WebsiteURL       string
	VersionsListLink string
	Tr               []string
	List             []resultslist
	headS
}

type headS struct {
	Description string
	Title       string
}

type pageS struct {
	headS
	Name       string
	List       []resultslist
	WebsiteURL string
	Link       string
}

type resultslist struct {
	Title  string
	Link   string
	Txt    template.HTML
	TdList []template.HTML
}

func (h *historyS) parse(w http.ResponseWriter) {
	h.Title += " - CurseForge mod"
	err := t.ExecuteTemplate(w, "history", h)
	if err != nil {
		log.Println(err)
	}
}

func (p *pageS) parse(w http.ResponseWriter, nextlink string) {
	if len(p.List) == 20 || nextlink != "" {
		p.Link = nextlink
	}
	p.Title += " - CurseForge mod"
	err := t.ExecuteTemplate(w, "page", p)
	if err != nil {
		log.Println(err)
	}
}

var t *template.Template

func init() {
	var err error
	t, err = template.ParseFS(htmlfs, "html/*")
	if err != nil {
		panic(err)
	}
	w := &bytes.Buffer{}
	type Title struct {
		Title       string
		Description string
	}
	err = t.ExecuteTemplate(w, "index", Title{Title: "CurseForge 搜索 - 搜索 CurseForge 上的东西并下载", Description: "搜索 CurseForge 上的东西并下载。"})
	if err != nil {
		panic(err)
	}
	index = w.Bytes()
}
