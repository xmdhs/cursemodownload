package curseapi

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

type cache struct {
	f       *fastcache.Cache
	cancel  func()
	expdate time.Duration
	bfpool  sync.Pool
}

func newcache(expdata time.Duration) *cache {
	c := &cache{}
	c.f = fastcache.New(32000000)
	c.expdate = expdata
	c.bfpool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
	return c
}

func (c *cache) Close() {
	c.cancel()
}

func (c *cache) Load(key string) []byte {
	b := c.f.GetBig(nil, []byte(key))
	if b == nil {
		return nil
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
	w := c.bfpool.Get().(*bytes.Buffer)
	binary.Write(w, binary.BigEndian, time.Now().Add(c.expdate).Unix())
	w.Write(adate)
	b := w.Bytes()
	c.f.SetBig([]byte(key), b)
	w.Reset()
	c.bfpool.Put(w)
}

func (c *cache) Delete(key string) {
	c.f.Del([]byte(key))
}
