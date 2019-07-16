package lru

import (
	"container/list"
	"sync"
)

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
	//If true, we move the entry to the front bcos it has been accessed
	if element, ok := c.cache[key]; ok {
		c.dll.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}
	//Else we just add the new entry element to the cache
	e := c.dll.PushFront(&entry{key: key, value: value})
	c.cache[key] = e

	//Remove the least used, if we exceed the capacity and since capacity is not enforced on the linked list
	if c.dll.Len() > c.capacity && c.capacity > 0 {
		c.removeLeastUsed()
	}
}

func (c *Cache) removeLeastUsed() (string, interface{}) {
	el := c.dll.Back()
	//List is empty
	if el == nil {
		return "", el
	}
	c.dll.Remove(el)
	value := el.Value.(*entry)
	delete(c.cache, value.key)
	return value.key, value.value
}

func (c *Cache) Get() {

}
