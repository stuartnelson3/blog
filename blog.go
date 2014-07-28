package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/pat"
)

func main() {
	var (
		mux        = pat.New()
		port       = flag.String("p", "3000", "address to bind the server on")
		dev        = flag.Bool("dev", false, "enable/disable dev mode")
		posts      = allPosts()
		show, _    = template.ParseFiles("./templates/show.html")
		index, _   = template.ParseFiles("./templates/index.html")
		newPost, _ = template.ParseFiles("./templates/new.html")
	)
	flag.Parse()

	if *dev {
		mux.Get("/new_post", func(w http.ResponseWriter, r *http.Request) {
			newPost.Execute(w, nil)
		})
		mux.Post("/new_post", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.FormValue("title"))
			p, err := CreatePost(r.FormValue("title"), r.FormValue("body"))
			if err != nil {
				http.Error(w, "Error making post.", 500)
				return
			}
			posts = allPosts()
			show.Execute(w, p)
		})
	}

	mux.Get("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	mux.Get("/{slug}", func(w http.ResponseWriter, r *http.Request) {
		p, err := findPost(r.URL.Query().Get(":slug"))
		if err != nil {
			http.Error(w, err.Error(), 404)
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
