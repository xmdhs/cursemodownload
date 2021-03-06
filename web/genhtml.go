package web

import (
	"html/template"
	"io"
	"log"
)

type results struct {
	Name       string
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
	t, err = template.New("page").Parse(searchhtml)
	if err != nil {
		panic(err)
	}
	t, err = t.New("body").Parse(body)
	if err != nil {
		panic(err)
	}
}
