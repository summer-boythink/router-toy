## router-toy

* A route with only 100 lines,
Facilitate personal study

* Fuzzy route matching is supported

## Usage
```go
package main

import (
	router "github.com/summer-boythink/router-toy"
	"net/http"
)

func main(){
    r := router.New()
    r.Get("aa/:name/c", func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte("aaa"))
    })
    r.Get("/q", func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte(("bbb")))
    })
    r.Post("/q/*all", func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte(("ccc")))
    })
    
    http.ListenAndServe("localhost:8080", r)
}
```
