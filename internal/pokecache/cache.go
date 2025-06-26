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
	fmt.Println("CACHE: Creating new cache with interval: ", interval)
	c := Cache{
		cacheData: make(map[string]cacheEntry),
		mu:        &sync.Mutex{},
		interval:  interval,
	}
	// could possibly pass the interval to the reapLoop so that
	// it will be aware of the cache delete interval and then
	// interval could be removed from the Cache struct
	cacheTicker := time.NewTicker(interval)
	go c.reapLoop(cacheTicker)

	return c
}

func (c Cache) Add(key string, val []byte) {
	fmt.Println("CACHE: Adding item to cache with key: ", key)
	fmt.Println("CACHE: len(val): ", len(val))
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheData[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	fmt.Println("CACHE: Added item to cache with key: ", key)
	fmt.Println("CACHE: Added item to cache with val: ", val)
}

func (c Cache) Get(key string) ([]byte, bool) {
	fmt.Println("CACHE: Looking in cache for item with key: ", key)
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cacheData[key]
	if !ok {
		fmt.Println("CACHE: Did not find item in cache with key: ", key)
		return []byte{}, false
	}

	fmt.Println("CACHE: Found item in cache with key: ", key)
	fmt.Println("CACHE: Returning cache data: ", val.val)
	return val.val, true
}

// blocks on the passed ticker channel until time data
// is sent at the interval specified at Cache creation
func (c Cache) reapLoop(reapTicker *time.Ticker) {

	for ; true; <-reapTicker.C {
		fmt.Println("CACHE: reapLoop is executing, time: ", time.Now())
		c.mu.Lock()
		for key, entry := range c.cacheData {
			timeSinceCreation := time.Since(entry.createdAt)
			if timeSinceCreation > c.interval {
				fmt.Println("CACHE: Deleting from cache item with key: ", key)
				delete(c.cacheData, key)
			}
		}
		c.mu.Unlock()
	}
}
