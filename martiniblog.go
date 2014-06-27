package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

type Post struct {
	Title     string        `json:"title"`
	Body      template.HTML `json:"body"`
	Slug      string        `json:"slug"`
	CreatedAt string        `json:"createdAt"`
	Mtime     int64         `json:"mtime"`
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout:     "layout",
		Extensions: []string{".html"}}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", Post{}.All())
	})

	m.Get("/:post", func(params martini.Params, r render.Render) {
		r.HTML(200, "show", Post{}.Find(params["post"]))
	})

	m.Run()
}

type ByMtime []*Post

func (a ByMtime) Len() int      { return len(a) }
func (a ByMtime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMtime) Less(i, j int) bool {
	return a[i].Mtime > a[j].Mtime
}

func (p Post) All() []*Post {
	matches, _ := filepath.Glob("posts/*.json")
	posts := make([]*Post, len(matches))
	for i, match := range matches {
		var post = &Post{}
		data, _ := ioutil.ReadFile(match)
		json.Unmarshal(data, post)
		posts[i] = post
	}
	sort.Sort(ByMtime(posts))
	return posts
}

func (p Post) Find(slug string) *Post {
	var post = &Post{}
	absPath, _ := filepath.Abs("posts/" + slug + ".json")
	data, _ := ioutil.ReadFile(absPath)
	json.Unmarshal(data, post)
	return post
}
