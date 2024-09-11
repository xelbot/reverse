# reverse
Golang URL reverse

Simple URL reverse package. It's useful for templates. You can get a URL by a name and params and not depend on URL structure.

It fits to any router. All it does is just stores urls by a name and replace params when you retrieve a URL.
To use it you have to add a URL with a name, raw URL with placeholders (params).

```go
// To set a URL and return raw URL use:
reverse.Add("UrlName", "/url_path/:param1/:param2")
// OUT: "/url_path/:param1/:param2"

// To set a URL with group (subrouter) prefix and return URL without prefix use:
reverse.Group("GroupName", "/prefix")
reverse.AddGr("UrlName", "GroupName", "/:param1/:param2")
// OUT: "/:param1/:param2"

// Note, that these funcs panic if errors.

// To retrieve a URL by name with given params use:
url, err := reverse.Get("UrlName", "param1", "value1", "param2", "value2")
// OUT: "/url_path/value1/value2"

// Get all url as map[string]string
reverse.GetAllUrls()
```

Example for Gin router (https://github.com/gin-gonic/gin):

```go
func main() {
    router := gin.Default()

    // URL: "/"
    // To fetch the url use: reverse.Get("home")
    router.GET(reverse.Add("home", "/"), indexEndpoint)

    // URL: "/get/123"
    // With param: c.Param("id")
    // To fetch the URL use: reverse.Get("get_url", "123")
    router.GET(reverse.Add("get_url", "/get/:id"), getUrlEndpoint)

    // Simple group: v1 (each URL starts with /v1 prefix)
    groupName := "v1"
    v1 := router.Group(reverse.Group(groupName, "/v1"))
    {
        // URL: "/v1"
        // To fetch the URL use: reverse.Get("v1_root")
        v1.GET(reverse.AddGr("v1_root", groupName, ""), v1RootEndpoint)

        // URL: "v1/read/cat123/id456"
        // With params (c.Param): catId, articleId
        // To fetch the URL use: reverse.Get("v1_read", "123", "456")
        v1.GET(reverse.AddGr("v1_read", groupName, "/read/cat:catId/id:articleId"), readEndpoint)

        // URL: /v1/login
        // To fetch the URL use: reverse.Get("v1_login")
        v1.GET(reverse.AddGr("v1_login", groupName, "/login"), loginGetEndpoint)
    }

    router.Run(":8080")
}

```

Example using Goji router:

```go
package main

import (
        "fmt"
        "net/http"
        "github.com/xelbot/reverse"
        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
        // We can get reversed URL by it's name and a list of params:
        // reverse.Get("UrlName", "value1", "value2")

        fmt.Fprintf(w, "Hello, %s", reverse.MustGet("HelloUrl", "name", c.URLParams["name"]))
}

func main() {
        // Set a URL and Params and return raw URL to a router
        // reverse.Add("UrlName", "/url_path/:param1/:param2", ":param1", ":param2")

        goji.Get(reverse.Add("HelloUrl", "/hello/:name"), hello)

        // In regexp instead of: re := regexp.MustCompile("^/comment/(?P<id>\\d+)$")
        re := regexp.MustCompile(reverse.Add("DeleteCommentUrl", "^/comment/(?P<id>\\d+)$"))
        goji.Delete(re, deleteComment)

        goji.Serve()
}
```

Example for Gorilla Mux

```go
// Original set: r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
r.HandleFunc(reverse.Add("ArticleCatUrl", "/articles/{category}/{id:[0-9]+}", "{category}", "{id:[0-9]+}"), ArticleHandler)

// So, if we want to retrieve URL "/articles/news/123", we call:
fmt.Println( reverse.MustGet("ArticleCatUrl", "category", "news", "id", "123") )
```

Example subrouters for [Chi](https://github.com/go-chi/chi) router:

```go
// Original code
r.Route("/articles", func(r chi.Router) {
	r.Get("/", listArticles)
	r.Route("/{articleID}", func(r chi.Router) {
		r.Get("/", getArticle)
	})
})

// With reverse package
r.Route(reverse.Group("articles", "/articles"), func(r chi.Router) {
	r.Get(reverse.AddGr("list_articles", "articles", "/"), listArticles)
	r.Route(reverse.Group("article", "/articles/{articleID}"), func(r chi.Router) {
		r.Get(reverse.AddGr("articles", "article", "/"), getArticle)
	})
})

// Get a reverse URL:
reverse.Get("get_article", "articleID", "123")
// Output: /articles/123/

// One more example (without tailing slashes)
r.Route(reverse.Group("admin", "/admin"), func(r chi.Router) {
	r.Get(reverse.AddGr("admin.index", "admin", "/"), index)

	r.Route(reverse.Group("admin.login", "/admin/login"), func(r chi.Router) {
		r.Get(reverse.AddGr("admin.login", "admin.login", "/"), login)
		r.Post("/", loginPost)
	})
})
```
