package itemdata

import (
	"ccfactory/server/debug"
	"encoding/gob"
	"os"
	"sync"
)

type DetailCache struct {
	path  string
	data  map[string]*Detail
	mutex sync.RWMutex
}

func NewDetailCache(path string) *DetailCache {
	c := &DetailCache{
		path: path,
		data: map[string]*Detail{},
	}

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		_ = gob.NewDecoder(file).Decode(&c.data)
	}

	return c
}

func (c *DetailCache) save() {
	defer debug.Timer("DetailCahe.save")()

	file, err := os.Create(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_ = gob.NewEncoder(file).Encode(c.data)
}

func (c *DetailCache) Get(id string) *Detail {
	defer debug.Timer("DetailCahe.Get")()

	c.mutex.RLock()
	data, ok := c.data[id]
	c.mutex.RUnlock()
	if ok {
		return data
	}
	return nil
}

func (c *DetailCache) GetOrSet(id string, fn func() *Detail) *Detail {
	defer debug.Timer("DetailCahe.GetOrSet")()

	c.mutex.RLock()
	data, ok := c.data[id]
	c.mutex.RUnlock()
	if ok {
		return data
	}

	data = fn()

	c.mutex.Lock()
	c.data[id] = data
	c.mutex.Unlock()

	c.save()

	return data
}
