package cache

import(
	"time"
	"sync"
)

// Entry is an entry in a cache
type Entry struct {
	Content []byte
	Expiration int64
}

// Cache is a simple cache, with get and set methods
type Cache struct {
	items map[string]Entry
	lock *sync.Mutex
}

// NewCache is the constructor for Cache
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]Entry),
		lock: &sync.Mutex{},
	}
	go c.monitorCache(10 * time.Second)
	return c
}

// Get returns the content from the Cache for the given key or nil 
//   if there is no entry for the key
func (cache *Cache) Get(key string) []byte  {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if entry, present := cache.items[key]; present  {
		return entry.Content
	}
	return nil
}

// monitorCache monitors the cache once per the given duration
func (cache *Cache) monitorCache(duration time.Duration)  {
	for {
		time.Sleep(duration)
		cache.lock.Lock()
		for key, entry := range cache.items  {
			if entry.Expiration > time.Now().UnixNano()  {
				delete(cache.items, key)
			}
		}
		cache.lock.Unlock()
	}
}

// Set sets the content for the entry at the given key to the given content, 
//   and keeps it for the given duration
func (cache *Cache) Set(key string, content []byte, duration time.Duration)  {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.items[key] = Entry{
		Content: content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}











