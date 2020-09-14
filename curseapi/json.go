package curseapi

import (
	"encoding/json"
	"fmt"
)

type Modinfo struct {
	Name                   string             `json:"name"`
	ID                     int                `json:"id"`
	GameVersionLatestFiles []gameVersionFiles `json:"gameVersionLatestFiles"`
}

type gameVersionFiles struct {
	GameVersion     string `json:"gameVersion"`
	ProjectFileId   int    `json:"projectFileId"`
	ProjectFileName string `json:"projectFileName"`
}

func json2Modinfo(jsonbyte []byte) ([]Modinfo, error) {
	m := make([]Modinfo, 0, 5)
	err := json.Unmarshal(jsonbyte, &m)
	if err != nil {
		return nil, fmt.Errorf("json2Modinfo: %w", err)
	}
	return m, nil
}
