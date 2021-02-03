package curseapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/groupcache"
)

var c = http.Client{Timeout: 10 * time.Second}

func httpget(url string) ([]byte, error) {
	reqs, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	b, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	return b, err
}

func httpcache(url string) ([]byte, error) {
	t := time.Now()
	ti := t.Round(time.Hour).Unix()
	k := key{
		URL:  url,
		Time: ti,
	}
	b, err := json.Marshal(k)
	if err != nil {
		return nil, fmt.Errorf("httpcache: %w", err)
	}
	bb := make([]byte, 0)
	err = cache.Get(context.TODO(), string(b), groupcache.AllocatingByteSliceSink(&bb))
	if err != nil {
		return nil, fmt.Errorf("httpcache: %w", err)
	}
	return bb, nil
}
