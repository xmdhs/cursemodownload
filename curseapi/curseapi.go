package curseapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

//From https://gaz492.github.io/TwitchAPI/

func Searchmod(key string, index string, sectionId int) ([]Modinfo, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/search?categoryId=0&gameId=432&index=` + index + `&pageSize=20&searchFilter=` + url.QueryEscape(key) + `&sectionId=` + strconv.Itoa(sectionId) + `&sort=0`
	b, err := httpcache(aurl)
	if err != nil {
		return nil, fmt.Errorf("Searchmod: %w", err)
	}
	m, err := json2Modinfo(b)
	if err != nil {
		return nil, fmt.Errorf("Searchmod: %w", err)
	}
	return m, nil
}

func FileId2downloadlink(id string) (string, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/0/file/` + id + `/download-url`
	b, err := httpcache(aurl)
	if err != nil {
		return "", fmt.Errorf("FileId2downloadlink: %w", err)
	}
	return string(b), nil
}

//https://media.forgecdn.net/files/3046/220/jei-1.16.2-7.3.2.25.jar

func AddonInfo(addonID string) (Modinfo, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/` + addonID
	b, err := httpcache(aurl)
	if err != nil {
		return Modinfo{}, fmt.Errorf("AddonInfo: %w", err)
	}
	m := Modinfo{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return Modinfo{}, fmt.Errorf("AddonInfo: %w", err)
	}
	sort.Slice(m.GameVersionLatestFiles, func(i, j int) bool {
		return m.GameVersionLatestFiles[i].ProjectFileId > m.GameVersionLatestFiles[j].ProjectFileId
	})
	return m, nil
}

func Addonfiles(addonID string) ([]Files, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/` + addonID + `/files`
	b, err := httpcache(aurl)
	if err != nil {
		return nil, fmt.Errorf("Addonfiles: %w", err)
	}
	m := make([]Files, 0)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("Addonfiles: %w", err)
	}
	sort.Slice(m, func(i, j int) bool {
		return m[i].ID > m[j].ID
	})
	return m, nil
}
