package curseapi

import (
	"fmt"
	"net/url"
)

//From https://gist.github.com/crapStone/9a423f7e97e64a301e88a2f6a0f3e4d9#file-curse_api-md

func Searchmod(key string, index string) ([]Modinfo, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/search?categoryId=0&gameId=432&index=` + index + `&pageSize=20&searchFilter=` + url.QueryEscape(key) + `&sectionId=6&sort=0`
	b, err := httpget(aurl)
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
	b, err := httpget(aurl)
	if err != nil {
		return "", fmt.Errorf("FileId2downloadlink: %w", err)
	}
	return string(b), nil
}
