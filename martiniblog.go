package main

import (
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "sort"
)

type Post struct {
    Title     string `json:"title"`
    Body      string `json:"body"`
    Slug      string `json:"slug"`
    CreatedAt string `json:"createdAt"`
    Mtime     int64  `json:"mtime"`
}

func main() {
    m := martini.Classic()
    m.Use(render.Renderer(render.Options{
        Layout:     "layout",
        Extensions: []string{".html"}}))

    m.Get("/", func(r render.Render) {
        posts := Post{}.All()
        r.HTML(200, "index", posts)
    })

    m.Get("/:post", func(params martini.Params, r render.Render) {
        post := Post{}.Find(params["post"])
        r.HTML(200, "show", post)
    })

    m.Run()
}

// sort posts by mtime
type ByMtime []*Post

func (a ByMtime) Len() int      { return len(a) }
func (a ByMtime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMtime) Less(i, j int) bool {
    return a[i].Mtime > a[j].Mtime
}

func (p Post) All() []*Post {
    var posts []*Post
    matches, _ := filepath.Glob("posts/*.json")
    for i := 0; i < len(matches); i++ {
        var post = &Post{}
        data, _ := ioutil.ReadFile(matches[i])
        json.Unmarshal(data, post)
        posts = append(posts, post)
    }
    sort.Sort(ByMtime(posts))
    return posts
}

func (p Post) Find(cond string) *Post {
    var post = &Post{}
    absPath, _ := filepath.Abs("posts/" + cond + ".json")
    data, _ := ioutil.ReadFile(absPath)
    json.Unmarshal(data, post)
    return post
}
