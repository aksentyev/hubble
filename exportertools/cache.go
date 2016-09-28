package exportertools

import (
    "github.com/patrickmn/go-cache"
    "time"
    "sync"
    "sort"
    "errors"
    "fmt"
)

// Cache for metrics values. Scraping metrics is async.
type Cache struct {
    *cache.Cache
    TTL time.Duration
    mu  *sync.Mutex
}

func NewCache(interval int) *Cache {
    ttl := time.Duration(interval) * time.Second
    cacheCleanUpInterval := time.Duration(interval / 2) * time.Second
    c := Cache{
        cache.New(ttl, cacheCleanUpInterval),
        ttl,
        new(sync.Mutex),
    }

    return &c
}

func (c *Cache) Set(m *Metric) {
    c.mu.Lock()
    c.Cache.Set(m.Name, m, c.TTL)
    c.mu.Unlock()
}

// GetActual returns actual values from the cache
func (c *Cache) Get(k string) (resp *Metric, err error) {
    c.mu.Lock()
    i, found := c.Cache.Get(k)
    c.mu.Unlock()
    if found {
        resp = i.(*Metric)
    } else {
        err = errors.New(fmt.Sprint("Cache has no value with key: %v", k))
    }
    return resp, err
}


func  (c *Cache) MetricNames() []string {
    c.mu.Lock()
    names := make([]string, 0)
    for k, _ := range c.Cache.Items() {
        names = append(names, k)
    }
    sort.Strings(names)
    c.mu.Unlock()
    return names
}
