package lru

import (
	"container/list"
	"sync"
)

//Cache Represents our LRU Cache
type Cache struct {
	mu       sync.Mutex
	cache    map[string]*list.Element
	dll      *list.List
	capacity int
}

//entry for list.Element Value interface field
type entry struct {
	key   string
	value interface{}
}

//New Instantiates a Cache
func New(maxSize int) *Cache {
	return &Cache{
		capacity: maxSize,
		cache:    make(map[string]*list.Element),
		dll:      list.New(),
	}
}

//Add - Adds  a new entry to the cache
func (c *Cache) Add(key string, value interface{}) {
	//Before adding, acquire lock
	c.mu.Lock()
	defer c.mu.Unlock()

	//Check if the entry is already in the cache
	if element, ok := c.cache[key]; ok {
		c.dll.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}

	e := c.dll.PushFront(&entry{key: key, value: value})
	c.cache[key] = e

	if c.dll.Len() > c.capacity && c.capacity > 0 {
		c.removeLeastUsed()
	}
}

func (c *Cache) removeLeastUsed() (string, interface{}) {
	el := c.dll.Back()
	if el == nil {
		return "", el
	}
	c.dll.Remove(el)
	value := el.Value.(*entry)
	delete(c.cache, value.key)
	return value.key, value.value
}

//Get - Retrieves a key from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.contains(key); ok {
		c.dll.MoveToFront(element)
		return element.Value.(*entry), ok
	}
	return nil, false
}

//Remove Removes an entry
func (c *Cache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.cache[key]; ok {
		c.dll.Remove(el)
		val := el.Value.(*entry)
		delete(c.cache, val.key)
		return true
	}
	return false
}

//Len Returns the length of the items in the cache
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.dll.Len()
}

//contains - Checks for the availability of a key in cache
//Returns the element and the boolean signifying availability
func (c *Cache) contains(key string) (*list.Element, bool) {
	entry, ok := c.cache[key]
	return entry, ok

}

//Contains - Checks for the availability of a key in cache
func (c *Cache) Contains(key string) bool {
	_, ok := c.cache[key]
	return ok
}
