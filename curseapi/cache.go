package curseapi

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

type cache struct {
	f       *fastcache.Cache
	cancel  func()
	expdate time.Duration
}

func newcache(expdata time.Duration) *cache {
	c := &cache{}
	c.f = fastcache.New(32000000)
	c.expdate = expdata
	return c
}

func (c *cache) Close() {
	c.cancel()
}

func (c *cache) Load(key string) []byte {
	b := c.f.GetBig(nil, []byte(key))
	if b == nil {
		b = c.f.Get(nil, []byte(key))
		if b == nil {
			return nil
		}
	}
	var d int64
	err := binary.Read(bytes.NewReader(b[:8]), binary.BigEndian, &d)
	if err != nil {
		return nil
	}
	t := time.Unix(d, 0)
	if t.Before(time.Now()) {
		c.f.Del([]byte(key))
		return nil
	}
	return b[8:]
}

func (c *cache) Store(key string, adate []byte) {
	w := bytes.NewBuffer(nil)
	binary.Write(w, binary.BigEndian, time.Now().Add(c.expdate).Unix())
	w.Write(adate)
	b := w.Bytes()

	if len(b) > 64000 {
		c.f.SetBig([]byte(key), b)
	} else {
		c.f.Set([]byte(key), b)
	}
}

func (c *cache) Delete(key string) {
	c.f.Del([]byte(key))
}
