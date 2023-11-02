package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu sync.Mutex 
	v map[string]string
}

func (c *Cache) LoadOrStore(key string, value string) (val string, ok bool) {
	c.mu.Lock()
	val, ok = c.v[key]
	if !ok {
		c.v[key] = value
	}
	c.mu.Unlock()
	return val, ok
}

func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	_, ok := c.v[key]
	delete(c.v, key)
	c.mu.Unlock()
	return ok 
}

var cache Cache

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, done chan bool) {
	if depth <= 0 {
		done <- true
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		done <- true
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	done_childs := make(chan bool)
    // i don't think we need to make a new channel
	childs := 0
	for _, u := range urls {
        // notice that we need load and store in a single atomic operation
		_, found := cache.LoadOrStore(u, "")
		if !found {
			childs += 1
			go Crawl(u, depth-1, fetcher, done_childs)
		}
	}
	for i := 0; i < childs; i++ {
		<-done_childs
	}
	done <- true
	return
}

//func main() {
//	cache = Cache{v: make(map[string]bool)}
//	cache.LoadOrStore("https://golang.org/", true)
//	done := make(chan bool)
//	go Crawl("https://golang.org/", 4, fetcher, done)
//	<- done
//}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

