# reverse
Go (golang) url reverse


Example using Goji router:

```go
package main

import (
        "fmt"
        "net/http"
        "github.com/alehano/reverse"
        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
        // We can get reversed url by it's name and a list of params:
        // reverse.Urls.MustReverse("UrlName", "param1", "param2")

        fmt.Fprintf(w, "Hello, %s", reverse.Urls.MustReverse("HelloUrl", c.URLParams["name"]))
}

func main() {
        // Set an Url and Params and return raw url to a router
        // reverse.Urls.MustAdd("UrlName", "/url_path/:param1/:param2", ":param1", ":param2")

        goji.Get(reverse.Urls.MustAdd("HelloUrl", "/hello/:name", ":name"), hello)
        
        // For regexp you can save url separately and then get it as usual:
        // reverse.Urls.MustReverse("deleteCommentUrl", "123")
        reverse.Urls.MustAdd("deleteCommentUrl", "/comment/:id", ":id")
        re := regexp.MustCompile("^/comment/(?P<id>\d+)$")
        goji.Delete(re, deleteComment)
        
        goji.Serve()
}
```
