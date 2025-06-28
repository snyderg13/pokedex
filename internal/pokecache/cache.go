package pokecache

import (
	"fmt"
	"sync"
	"time"
)

var cacheDebug bool = false
var cacheDelDebug bool = true

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheData map[string]cacheEntry
	mu        *sync.Mutex
}

// creates new Cache and launches the reapLoop as a go routine
func NewCache(interval time.Duration) Cache {
	if cacheDebug {
		fmt.Println("CACHE: Creating new cache with interval: ", interval)
	}

	c := Cache{
		cacheData: make(map[string]cacheEntry),
		mu:        &sync.Mutex{},
	}
	go c.reapLoop(interval)

	return c
}

func (c Cache) Add(key string, val []byte) {
	if cacheDebug {
		fmt.Println("CACHE: Adding item to cache with key: ", key)
		fmt.Println("CACHE: len(val): ", len(val))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheData[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	if cacheDebug {
		fmt.Println("CACHE: Added item to cache with key: ", key)
		fmt.Println("CACHE: Added item to cache with val: ", val)
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	if cacheDebug {
		fmt.Println("CACHE: Looking in cache for item with key: ", key)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cacheData[key]
	if !ok {
		if cacheDebug {
			fmt.Println("CACHE: Did not find item in cache with key: ", key)
		}

		return []byte{}, false
	}

	if cacheDebug {
		fmt.Println("CACHE: Found item in cache with key: ", key)
		fmt.Println("CACHE: Returning cache data: ", val.val)
	}

	return val.val, true
}

// blocks on the passed ticker channel until time data
// is sent at the interval specified at Cache creation
func (c Cache) reapLoop(interval time.Duration) {
	reapTicker := time.NewTicker(interval)
	for ; true; <-reapTicker.C {
		if cacheDebug {
			fmt.Println("CACHE: reapLoop is executing, time: ", time.Now())
		}
		c.mu.Lock()
		for key, entry := range c.cacheData {
			if time.Since(entry.createdAt) > interval {
				if cacheDebug || cacheDelDebug {
					fmt.Println("CACHE: Deleting from cache item with key: ", key)
				}
				delete(c.cacheData, key)
			}
		}
		c.mu.Unlock()
	}
}
