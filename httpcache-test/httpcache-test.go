package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gregjones/httpcache"
)

func main() {
	// shortcut, but we go the long way to test reading the cache
	// client := http.Client{Transport: httpcache.NewMemoryCacheTransport()}

	cache := httpcache.NewMemoryCache()
	transport := httpcache.NewTransport(cache)
	client := http.Client{Transport: transport}
	response, err := client.Get("http://www.example.com")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// cache is only set when we read the body
	buf := bytes.Buffer{}
	buf.ReadFrom(response.Body)
	buf.String()

	request, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		panic(err)
	}
	response, err = httpcache.CachedResponse(cache, request)
	if err != nil {
		panic(err)
	}

	if response == nil {
		panic("response is nil")
	}
	os.Exit(0)
}
