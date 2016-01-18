package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	bf "github.com/russross/blackfriday"
)

type Post struct {
	Title     string        `json:"title"`
	Slug      string        `json:"slug"`
	CreatedAt string        `json:"createdAt"`
	Body      template.HTML `json:"body"`
	Mtime     int64         `json:"mtime"`
}

func allPosts() []*Post {
	// replace this with db stuff
	paths, _ := filepath.Glob("posts/*.json")
	posts := make([]*Post, len(paths))
	for i, path := range paths {
		post := &Post{}
		f, err := os.Open(path)
		if err != nil {
			log.Printf("error opening file: %v", err)
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
	// replace this with db stuff
	f, err := os.Open("./posts/" + slug + ".json")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	post := &Post{}
	json.NewDecoder(f).Decode(post)
	return post, nil
}

type ByMtime []*Post

func (a ByMtime) Len() int           { return len(a) }
func (a ByMtime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMtime) Less(i, j int) bool { return a[i].Mtime > a[j].Mtime }

func CreatePost(title, body string) (*Post, error) {
	r := bf.HtmlRenderer(
		bf.HTML_HREF_TARGET_BLANK,
		"",
		"",
	)
	ext := bf.EXTENSION_NO_INTRA_EMPHASIS | bf.EXTENSION_FENCED_CODE | bf.EXTENSION_STRIKETHROUGH | bf.EXTENSION_LAX_HTML_BLOCKS
	p := &Post{
		Title:     title,
		Body:      template.HTML(string(bf.Markdown([]byte(body), r, ext))),
		Slug:      CreateSlug(title),
		CreatedAt: time.Now().Format("Jan 2 2006"),
		Mtime:     time.Now().Unix(),
	}

	// replace save json with a function to commit to db
	err := p.SaveJson()
	return p, err
}

func (p *Post) SaveJson() error {
	f, err := os.Create("./posts/" + p.Slug + ".json")
	if err != nil {
		return err
	}
	defer f.Close()

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
