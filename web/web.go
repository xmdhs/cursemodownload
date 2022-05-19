package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xmdhs/cursemodownload/curseapi"
)

var sectionIds = map[string]int{
	"1": 6,
	"2": 12,
	"3": 4471,
	"4": 17,
}

func WebRoot(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	page := req.FormValue("page")
	if page == "" {
		page = "0"
	}
	atype := req.FormValue("type")
	sectionId, ok := sectionIds[atype]
	if !ok {
		sectionId = 6
	}
	if len(query) > 100 {
		err := errors.New("关键词过长")
		http.Error(w, err.Error(), 500)
		return
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, err := search(query, page, sectionId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(r) == 0 {
		http.NotFound(w, req)
		return
	}
	i++
	page = strconv.FormatInt(i, 10)
	link := "./s?q=" + query + "&type=" + atype + "&page=" + page

	p := pageS{
		headS: headS{
			Description: "",
			Title:       query,
		},
		Name: query,
		List: r,
		Link: "",
	}
	p.parse(w, link)
}

var index []byte

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write(index)
}

func search(txt, offset string, sectionId int) ([]resultslist, error) {
	c, err := curseapi.Searchmod(txt, offset, sectionId)
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
	id := req.FormValue("id")
	c, err := curseapi.AddonInfo(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var r []resultslist
	var title string
	set := make(map[string]struct{})
	sort.Slice(c.GameVersionLatestFiles, func(i, j int) bool {
		v1 := c.GameVersionLatestFiles[i].GameVersion
		v2 := c.GameVersionLatestFiles[j].GameVersion
		v1l := strings.Split(v1, ".")
		v2l := strings.Split(v2, ".")
		a := 0
		if len(v1l) < len(v2l) {
			a = len(v1l)
		} else {
			a = len(v2l)
		}
		for i := 0; i < a; i++ {
			vn1, err := strconv.Atoi(v1l[i])
			if err != nil {
				vn1 = 0
			}
			vn2, err := strconv.Atoi(v2l[i])
			if err != nil {
				vn2 = 0
			}
			if vn1 > vn2 {
				return true
			} else if vn1 < vn2 {
				return false
			}
		}
		return len(v1l) > len(v2l)
	})

	vers := strings.Builder{}

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

				vers.WriteString(v.GameVersion + " ")

				r = append(r, temp)
			}
		}
	}

	p := pageS{
		headS: headS{
			Description: c.Name + " - " + c.Summary + " - files download " + vers.String(),
			Title:       title,
		},
		Name:       c.Name,
		List:       r,
		Link:       "",
		WebsiteURL: c.Links.WebsiteUrl,
	}
	p.parse(w, "")
}

func Getdownloadlink(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	id := q.Get("id")
	if id == "" {
		http.Error(w, "", 500)
		return
	}
	link, err := curseapi.FileId2downloadlink(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, req, link, http.StatusFound)
}

func History(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if len(q["id"]) == 0 || len(q["ver"]) == 0 {
		http.Error(w, "", 500)
		return
	}
	id, ver := q["id"][0], q["ver"][0]
	ch := make(chan curseapi.Modinfo, 10)
	errCh := make(chan error, 10)
	go func() {
		info, err := curseapi.AddonInfo(id)
		if err != nil {
			errCh <- err
		}
		ch <- info
	}()
	h, err := curseapi.Addonfiles(id, ver)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var info curseapi.Modinfo
	select {
	case info = <-ch:
	case err = <-errCh:
		http.Error(w, err.Error(), 500)
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
		tdlist := make([]template.HTML, 0)
		if v.DownloadUrl == "" {
			tdlist = append(tdlist, template.HTML(template.HTMLEscapeString(v.FileName+" (can not get download link)")))
		} else {
			tdlist = append(tdlist, template.HTML(`<a href="`+template.HTMLEscapeString(v.DownloadUrl)+`" target="_blank">`+template.HTMLEscapeString(v.FileName)+`</a>`))
		}
		tdlist = append(tdlist, template.HTML(template.HTMLEscapeString(releaseType[v.ReleaseType])))
		atime, err := time.Parse(time.RFC3339, v.FileDate)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		tdlist = append(tdlist, template.HTML(strconv.FormatInt(atime.Unix(), 10)))
		d := dependenciespase(v.Dependencies)
		if d == "" {
			d = "none"
		}
		tdlist = append(tdlist, template.HTML(d))
		r = append(r, resultslist{
			TdList: tdlist,
		})
	}
	hs := historyS{
		Name:             info.Name,
		Version:          ver,
		WebsiteURL:       info.Links.WebsiteUrl,
		VersionsListLink: "/curseforge/info?id=" + id,
		Tr:               tdname,
		List:             r,
		headS: headS{
			Description: info.Name + " - " + ver + " - " + info.Summary + " - files download",
			Title:       info.Name + " - " + ver,
		},
	}
	hs.parse(w)
}

var tdname = []string{"File Name", "Release Type", "File Date", "Dependencies"}

var releaseType = map[int]string{
	1: "Release",
	2: "Beta",
	3: "Alpha",
}

func dependenciespase(dependencies []curseapi.Dependencies) string {
	s := strings.Builder{}
	i := 0
	for _, v := range dependencies {
		if v.Type == 3 {
			s.WriteString(`<a href="` + dependencies2url(v) + `" target="_blank">` + strconv.Itoa(v.AddonId) + `</a> `)
			i++
		}
	}
	if i == 0 {
		return ""
	}
	return s.String()
}

func dependencies2url(dependencies curseapi.Dependencies) string {
	return "./info?id=" + strconv.Itoa(dependencies.AddonId)
}
