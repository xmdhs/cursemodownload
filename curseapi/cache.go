package curseapi

import (
	"context"
	"encoding/json"

	"github.com/golang/groupcache"
)

var cache *groupcache.Group

func init() {
	cache = groupcache.NewGroup("curseapi", 10000000, groupcache.GetterFunc(func(ctx context.Context, akey string, dest groupcache.Sink) error {
		k := key{}
		err := json.Unmarshal([]byte(akey), &k)
		if err != nil {
			return err
		}
		b, err := httpget(k.URL)
		if err != nil {
			return err
		}
		dest.SetBytes(b)
		return nil
	}))
}

type key struct {
	Time int64
	URL  string
}
