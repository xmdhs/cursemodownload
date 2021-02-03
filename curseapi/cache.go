package curseapi

import (
	"context"
	"sync"
	"time"
)

type cache struct {
	m      sync.Map
	cancel func()
}

type date struct {
	Time int64
	Date []byte
}

func newcache() *cache {
	c := &cache{}
	cxt := context.Background()
	cxt, cancel := context.WithCancel(cxt)
	c.cancel = cancel
	c.delete(cxt)
	return c
}

func (c *cache) delete(cxt context.Context) {
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for {
			c.m.Range(func(key, value interface{}) bool {
				d := value.(date)
				if time.Now().Unix()-d.Time > 1800 {
					c.m.Delete(key)
				}
				return true
			})

			select {
			case <-cxt.Done():
				return
			case <-t.C:
			}
		}
	}()
}

func (c *cache) Close() {
	c.cancel()
}

func (c *cache) Load(key string) []byte {
	t, ok := c.m.Load(key)
	if !ok {
		return nil
	}
	d, ok := t.(date)
	if !ok {
		return nil
	}
	return d.Date
}

func (c *cache) Store(key string, adate []byte) {
	d := date{
		Date: adate,
		Time: time.Now().Unix(),
	}
	c.m.Store(key, d)
}
