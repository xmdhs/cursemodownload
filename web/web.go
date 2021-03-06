package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/xmdhs/cursemodownload/curseapi"
)

func WebRoot(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	var page, query string
	if len(q["q"]) == 0 {
		query = ""
	} else {
		query = q["q"][0]
	}
	if len(q["page"]) == 0 {
		page = "0"
	} else {
		page = q["page"][0]
	}
	if len(query) > 100 {
		err := errors.New("关键词过长")
		e(w, err)
		return
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		e(w, err)
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, err := search(query, page)
	if err != nil {
		e(w, err)
		return
	}
	if len(r) == 0 {
		http.NotFound(w, req)
		return
	}
	i++
	page = strconv.FormatInt(i, 10)
	pase(w, r, query, page, "./s?q=", "")
}

func e(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(index))
}

func search(txt, offset string) ([]resultslist, error) {
	c, err := curseapi.Searchmod(txt, offset)
	if err != nil {
		return nil, err
	}
	r := make([]resultslist, 0, len(c))
	for _, v := range c {
		temp := resultslist{
			Title: v.Name,
			Link:  fmt.Sprintf("./info?id=%v", v.ID),
			Txt:   template.HTML(template.HTMLEscapeString(v.Summary)),
		}
		r = append(r, temp)
	}
	return r, nil
}

func Info(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if len(q["id"]) == 0 {
		e(w, errors.New(`""`))
		return
	}
	id := q["id"][0]
	c, err := curseapi.AddonInfo(id)
	if err != nil {
		e(w, err)
		return
	}
	var r []resultslist
	var title string
	set := make(map[string]struct{})
	if strconv.Itoa(c.ID) == id {
		title = c.Name
		r = make([]resultslist, 0, len(c.GameVersionLatestFiles))
		for _, v := range c.GameVersionLatestFiles {
			if _, ok := set[v.GameVersion]; !ok {
				set[v.GameVersion] = struct{}{}
				link := `./history?id=` + id + "&ver=" + v.GameVersion
				temp := resultslist{
					Title: template.HTMLEscapeString(v.GameVersion),
					Link:  link,
				}
				r = append(r, temp)
			}
		}
	}

	pase(w, r, title, "", "", c.WebsiteUrl)
}

func Getdownloadlink(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	id := q.Get("id")
	if id == "" {
		e(w, errors.New(`""`))
		return
	}
	link, err := curseapi.FileId2downloadlink(id)
	if err != nil {
		e(w, err)
		return
	}
	http.Redirect(w, req, link, 302)
}

func History(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if len(q["id"]) == 0 || len(q["ver"]) == 0 {
		e(w, errors.New(`""`))
		return
	}
	id, ver := q["id"][0], q["ver"][0]
	h, err := curseapi.Addonfiles(id)
	if err != nil {
		e(w, err)
		return
	}
	files := make([]curseapi.Files, 0)
	for _, v := range h {
		for _, vv := range v.GameVersion {
			if vv == ver {
				files = append(files, v)
			}
		}
	}
	r := make([]resultslist, 0)
	for _, v := range files {
		r = append(r, resultslist{
			Title: v.FileName + " " + releaseType[v.ReleaseType],
			Link:  v.DownloadUrl,
			Txt:   template.HTML(dependenciespase(v.Dependencies)),
		})
	}
	pase(w, r, id+" "+ver, "", "", "")
}

var releaseType = map[int]string{
	1: "Release",
	2: "Beta",
	3: "Alpha",
}

func dependenciespase(dependencies []curseapi.Dependencies) string {
	s := strings.Builder{}
	i := 0
	s.WriteString(`<p>依赖：`)
	for _, v := range dependencies {
		if v.Type == 3 {
			s.WriteString(`<a href="` + dependencies2url(v) + `" target="_blank">` + strconv.Itoa(v.AddonId) + `</a> `)
			i++
		}
	}
	s.WriteString(`</p>`)
	if i == 0 {
		return ""
	}
	return s.String()
}

func dependencies2url(dependencies curseapi.Dependencies) string {
	return "./info?id=" + strconv.Itoa(dependencies.AddonId)
}
