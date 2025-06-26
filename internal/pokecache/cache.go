package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheData map[string]cacheEntry
	mu        *sync.Mutex
	interval  time.Duration
}

// creates new Cache and launches the reapLoop as a go routine
func NewCache(interval time.Duration) Cache {
	c := Cache{
		cacheData: make(map[string]cacheEntry),
		mu:        &sync.Mutex{},
		interval:  interval,
	}
	cacheTicker := time.NewTicker(interval)
	go c.reapLoop(cacheTicker)

	return c
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheData[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cacheData[key]
	if !ok {
		return []byte{}, false
	}

	return val.val, true
}

// blocks on the passed ticker channel until time data
// is sent at the interval specified at Cache creation
func (c Cache) reapLoop(reapTicker *time.Ticker) {

	for ; true; <-reapTicker.C {
		fmt.Println("reapLoop is executing, time: ", time.Now())
		c.mu.Lock()
		for key, entry := range c.cacheData {
			timeSinceCreation := time.Since(entry.createdAt)
			if timeSinceCreation > c.interval {
				delete(c.cacheData, key)
			}
		}
		c.mu.Unlock()
	}
}
