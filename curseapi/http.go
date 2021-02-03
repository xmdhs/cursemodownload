package curseapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/groupcache/singleflight"
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

var acache = newcache()

var s = singleflight.Group{}

func httpcache(url string) ([]byte, error) {
	b := acache.Load(url)
	if b != nil {
		return b, nil
	}
	t, err := s.Do(url, func() (interface{}, error) {
		b, err := httpget(url)
		if err != nil {
			return nil, err
		}
		return b, nil
	})
	if err != nil {
		return nil, fmt.Errorf("httpcache: %w", err)
	}
	b = t.([]byte)
	acache.Store(url, b)
	return b, nil
}
