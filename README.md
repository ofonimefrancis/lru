#  lru
An LRU Cache Safe for concurrent access

This is a simple implementation of an LRU cache in Go. 
Concurrent access is achieved by acquiring a lock on the cache before addition or retrievals.

`Documentation`
==================
Documentation can be found on [Godoc](http://godoc.org/github.com/ofonimefrancis/lru)

`Usage`
================

```go
    l := lru.New(4)

    l.Add("one", 1)
    l.Add("two", 2)
    l.Add("three", 3)
    l.Add("four", 4)

    fmt.Println(l.Get("four"))

```


