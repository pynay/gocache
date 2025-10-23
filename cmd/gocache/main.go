package main

import (
	"fmt"
	"time"
	"github.com/pynay/gocache/internal/cache"
)


func main() {
    c := cache.New()
    stop := make(chan struct{})
    c.StartJanitor(2*time.Second, stop)

    c.Put("hello", []byte("world"), 3*time.Second)

    v, _ := c.Get("hello")
    fmt.Println("value:", string(v)) // prints "world"

    time.Sleep(4 * time.Second)
    v, err := c.Get("hello")
    fmt.Println("after 4s:", v, err) // should show expired
}