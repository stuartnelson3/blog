package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/gorilla/handlers"
	"github.com/gorilla/pat"
)

type Post struct {
	Title     string        `json:"title"`
	Body      template.HTML `json:"body"`
	Slug      string        `json:"slug"`
	CreatedAt string        `json:"createdAt"`
	Mtime     int64         `json:"mtime"`
}

func allPosts() []*Post {
	paths, _ := filepath.Glob("posts/*.json")
	posts := make([]*Post, len(paths))
	for i, path := range paths {
		post := &Post{}
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		json.NewDecoder(f).Decode(post)
		posts[i] = post
		f.Close()
	}
	sort.Sort(ByMtime(posts))
	return posts
}

func findPost(slug string) (*Post, error) {
	f, err := os.Open("./posts/" + slug + ".json")
	if err != nil {
		return nil, err
	}
	post := &Post{}
	json.NewDecoder(f).Decode(post)
	f.Close()

	return post, nil
}

type ByMtime []*Post

func (a ByMtime) Len() int           { return len(a) }
func (a ByMtime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMtime) Less(i, j int) bool { return a[i].Mtime > a[j].Mtime }

func main() {
	var (
		mux      = pat.New()
		port     = flag.String("p", "3000", "address to bind the server on")
		show, _  = template.ParseFiles("./templates/show.html")
		index, _ = template.ParseFiles("./templates/index.html")
		posts    = allPosts()
	)
	flag.Parse()

	mux.Get("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	mux.Get("/{slug}", func(w http.ResponseWriter, r *http.Request) {
		p, err := findPost(r.URL.Query().Get(":slug"))
		if err != nil {
			http.Error(w, "Post not found.", 404)
			return
		}
		show.Execute(w, p)
	})

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		index.Execute(w, posts)
	})

	log.Printf("Starting server on %s\n", *port)
	handler := handlers.LoggingHandler(os.Stdout, mux)
	handler = handlers.CompressHandler(handler)
	log.Fatal(http.ListenAndServe(":"+*port, handler))
}
