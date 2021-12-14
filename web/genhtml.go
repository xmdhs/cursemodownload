package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type results struct {
	Name        string
	Title       string
	List        []resultslist
	Link        string
	T           bool
	WebsiteURL  string
	Description string
	Table       bool
	Script      template.HTML
	Tr          []string
}

type resultslist struct {
	Title  string
	Link   string
	Txt    template.HTML
	TdList []template.HTML
}

func tablepase(w http.ResponseWriter, list []resultslist, Name, nextlink, titilelink, Description string, script template.HTML, tr []string) {
	T := true
	if len(list) != 20 || nextlink == "" {
		T = false
	}
	Table := false
	for _, v := range list {
		if len(v.TdList) > 0 {
			Table = true
			break
		}
	}
	r := results{
		Title:       Name + " - CurseForge mod",
		Name:        Name,
		Link:        nextlink,
		List:        list,
		WebsiteURL:  titilelink,
		T:           T,
		Description: Description,
		Table:       Table,
		Script:      script,
		Tr:          tr,
	}
	err := t.ExecuteTemplate(w, "page", r)
	if err != nil {
		log.Println(err)
		return
	}
}

func pase(w http.ResponseWriter, list []resultslist, Name, nextlink, titilelink, Description string) {
	tablepase(w, list, Name, nextlink, titilelink, Description, "", nil)
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
