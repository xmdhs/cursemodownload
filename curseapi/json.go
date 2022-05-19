package curseapi

import (
	"encoding/json"
	"fmt"
)

type apidata[v any] struct {
	Data v `json:"data"`
}

type Modinfo struct {
	Name                   string             `json:"name"`
	ID                     int                `json:"id"`
	GameVersionLatestFiles []GameVersionFiles `json:"latestFilesIndexes"`
	Links                  struct {
		WebsiteUrl string `json:"websiteUrl"`
	} `json:"links"`
	Summary string `json:"summary"`
}

type GameVersionFiles struct {
	GameVersion     string `json:"gameVersion"`
	ProjectFileId   int    `json:"fileId"`
	ProjectFileName string `json:"filename"`
}

func json2Modinfo(jsonbyte []byte) ([]Modinfo, error) {
	d := apidata[[]Modinfo]{}
	err := json.Unmarshal(jsonbyte, &d)
	if err != nil {
		return nil, fmt.Errorf("json2Modinfo: %w", err)
	}
	return d.Data, nil
}

type Files struct {
	ID           int            `json:"id"`
	FileName     string         `json:"fileName"`
	FileDate     string         `json:"fileDate"`
	Dependencies []Dependencies `json:"dependencies"`
	DownloadUrl  string         `json:"downloadUrl"`
	GameVersion  []string       `json:"gameVersions"`
	ReleaseType  int            `json:"releaseType"`
}

type Dependencies struct {
	AddonId int `json:"modId"`
	Type    int `json:"relationType"`
}
