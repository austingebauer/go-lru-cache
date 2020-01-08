# go-lru-cache [![GoReportCard](https://goreportcard.com/badge/github.com/austingebauer/go-lru-cache)](https://goreportcard.com/report/github.com/austingebauer/go-lru-cache) [![GoDoc](https://godoc.org/github.com/austingebauer/go-lru-cache?status.svg)](https://godoc.org/github.com/austingebauer/go-lru-cache)

A Go library that provides a thread-safe least recently used (LRU) cache with a fixed 
capacity.

## Installation

To install `go-lru-cache`, use `go get`.

```bash
go get github.com/austingebauer/go-lru-cache
```

Then import the library into your Go program.

```go
import "github.com/austingebauer/go-lru-cache"
```

## Usage

### API

`go-lru-cache` has a simple API.

It provides a `Put()` method that allows you to place key/value pairs into the cache.

It provides a `Get()` method that allows you to retrieve values given a key.

It provides a `Purge()` method that clears the cache.

It provides a `Len()` method that returns the length of the cache.

Please see the [GoDoc](https://godoc.org/github.com/austingebauer/go-lru-cache) for 
additional API documentation of the library.

**Example 1**: Basic usage
```go
cache, err := lru.NewCache(2, nil)

cache.Put(1, 2)
cache.Put(2, 3)
cache.Get(1)       // returns 2
cache.Put(3, 4)    // evicts 2->3
cache.Get(2)       // returns -1 (not found)
cache.Put(4, 5)    // evicts 1->2
cache.Get(1)       // returns -1 (not found)
cache.Get(3)       // returns 4
cache.Get(4)       // returns 5
```

**Example 2**: Usage showing `onEvict()` function passed into `NewCache()`
```go
keyVal := ""
onEvict := func(key, value interface{}) {
    keyVal = key.(string) + value.(string)
}

cache, err := NewCache(1, onEvict)
cache.Put("ada", "lovelace")
cache.Put("linus", "torvalds") // evicts ada->lovelace

// Eviction happened on ada->lovelace, thereby calling onEvict()
// So, printing keyVal will print "adalovelace"
fmt.Println(keyVal)
```

### Behavior

`go-lru-cache` will begin to evict the least recently used key/value pair when it has 
reached its given capacity.

Calls to both `Get()` and `Put()` count as usage of a given key/value pair.

If an `onEvicted` function has been passed into `NewCache()` during construction, then 
the function will be called with the evicted key/value pair every time an eviction occurs.

## Contributing

Pull requests are welcome. 

For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests along with changes.

## License

[MIT](LICENSE)
