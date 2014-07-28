package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	Body      string `json:"body"`
	Mtime     int64  `json:"mtime"`
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

func CreatePost(title, body string) (*Post, error) {
	p := &Post{
		Title:     title,
		Body:      string(blackfriday.MarkdownCommon([]byte(body))),
		Slug:      CreateSlug(title),
		CreatedAt: time.Now().Format("Jan 2 2006"),
		Mtime:     time.Now().Unix(),
	}

	err := p.SaveJson()
	return p, err
}

func (p *Post) SaveJson() error {
	f, err := os.Create("./posts/" + p.Slug + ".json")
	if err != nil {
		return err
	}
	j, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = f.Write(j)
	if err != nil {
		return err
	}
	return nil
}

func CreateSlug(title string) string {
	slug := strings.Join(strings.Fields(strings.ToLower(title)), "-")
	re := regexp.MustCompile("[^0-9A-Za-z_-]")
	return re.ReplaceAllLiteralString(slug, "")
}
