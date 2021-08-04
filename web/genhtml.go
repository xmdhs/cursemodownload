package web

import (
	"html/template"
	"io"
	"log"
)

type results struct {
	Name       string
	Title      string
	List       []resultslist
	Link       string
	T          bool
	WebsiteURL string
}

type resultslist struct {
	Title string
	Link  string
	Txt   template.HTML
}

func pase(w io.Writer, list []resultslist, Name, page, link, titilelink string) {
	T := true
	Link := ""
	if len(list) != 20 || page == "" {
		T = false
	} else {
		Link = link + Name + "&page=" + page
	}
	r := results{
		Title:      Name + " - curseforge mod",
		Name:       Name,
		Link:       Link,
		List:       list,
		WebsiteURL: titilelink,
		T:          T,
	}
	err := t.ExecuteTemplate(w, "page", r)
	if err != nil {
		log.Println(err)
		return
	}
}

var t *template.Template

func init() {
	var err error
	t, err = template.ParseFS(htmlfs, "html/*")
	if err != nil {
		panic(err)
	}
}
