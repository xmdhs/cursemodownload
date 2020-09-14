package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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
	pase(w, r, query, page, "./s?q=")
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
	if strconv.Itoa(c.ID) == id {
		title = c.Name
		r = make([]resultslist, 0, len(c.GameVersionLatestFiles))
		for _, v := range c.GameVersionLatestFiles {
			link := `./download?id=` + strconv.Itoa(v.ProjectFileId)
			temp := resultslist{
				Title: v.ProjectFileName + "  " + v.GameVersion,
				Link:  link,
				Txt:   template.HTML(`<a href="` + link + `" target="_blank">官方下载</a> <a href="` + link + "&cdn=1" + `" target="_blank">镜像下载</a>`),
			}
			r = append(r, temp)
		}
	}

	pase(w, r, title, "", "")
}

func Getdownloadlink(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if len(q["id"]) == 0 {
		e(w, errors.New(`""`))
		return
	}
	id := q["id"][0]
	link, err := curseapi.FileId2downloadlink(id)
	if err != nil {
		e(w, err)
		return
	}
	if len(q["cdn"]) != 0 {
		link = `https://cors.xmdhs.top/` + link
	}
	http.Redirect(w, req, link, 302)
}
