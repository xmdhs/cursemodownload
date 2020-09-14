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

func FileId2downloadlink(filename, id string) (string, error) {
	return `https://edge.forgecdn.net/files/` + id[:4] + "/" + id[4:] + "/" + filename, nil
}

//https://media.forgecdn.net/files/3046/220/jei-1.16.2-7.3.2.25.jar
