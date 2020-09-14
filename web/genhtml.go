package web

import (
	"html/template"
	"io"
	"log"
)

type results struct {
	Name string
	List []resultslist
	Link string
	T    bool
}

type resultslist struct {
	Title string
	Link  string
	Txt   template.HTML
}

func pase(w io.Writer, list []resultslist, Name, page, link string) {
	T := true
	Link := ""
	if len(list) != 20 || page == "" {
		T = false
	} else {
		Link = link + Name + "&page=" + page
	}
	r := results{
		Name: Name,
		Link: Link,
		List: list,
		T:    T,
	}
	t, err := template.New("page").Parse(searchhtml)
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(w, r)
	if err != nil {
		log.Println(err)
		return
	}
}
