package curseapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

//From https://gaz492.github.io/TwitchAPI/

const api = `https://api.curseforge.com/v1`

func Searchmod(key string, index string, sectionId int) ([]Modinfo, error) {
	aurl := api + `/mods/search?categoryId=0&gameId=432&index=` + index + `&pageSize=20&searchFilter=` + url.QueryEscape(key) + `&classId=` + strconv.Itoa(sectionId) + `&sortField=2&sortOrder=desc`
	b, err := httpcache(aurl, acache)
	if err != nil {
		return nil, fmt.Errorf("Searchmod: %w", err)
	}
	m, err := json2Modinfo(b)
	if err != nil {
		acache.Delete(aurl)
		return nil, fmt.Errorf("Searchmod: %w", err)
	}
	return m, nil
}

func FileId2downloadlink(id string) (string, error) {
	aurl := api + `/mods/0/file/` + id + `/download-url`
	b, err := httpcache(aurl, acache)
	if err != nil {
		return "", fmt.Errorf("FileId2downloadlink: %w", err)
	}
	return string(b), nil
}

//https://media.forgecdn.net/files/3046/220/jei-1.16.2-7.3.2.25.jar

func AddonInfo(addonID string) (Modinfo, error) {
	aurl := api + `/mods/` + addonID
	b, err := httpcache(aurl, acache)
	if err != nil {
		return Modinfo{}, fmt.Errorf("AddonInfo: %w", err)
	}
	d := apidata[Modinfo]{}
	err = json.Unmarshal(b, &d)
	m := d.Data
	if err != nil {
		acache.Delete(aurl)
		return Modinfo{}, fmt.Errorf("AddonInfo: %w", err)
	}
	return m, nil
}

func Addonfiles(addonID, gameVersion string) ([]Files, error) {
	aurl := api + `/mods/` + addonID + `/files?pageSize=10000&gameVersion=` + gameVersion
	b, err := httpcache(aurl, acache)
	if err != nil {
		return nil, fmt.Errorf("Addonfiles: %w", err)
	}
	d := apidata[[]Files]{}
	err = json.Unmarshal(b, &d)
	m := d.Data
	if err != nil {
		acache.Delete(aurl)
		return nil, fmt.Errorf("Addonfiles: %w", err)
	}
	sort.Slice(m, func(i, j int) bool {
		return m[i].ID > m[j].ID
	})
	return m, nil
}
