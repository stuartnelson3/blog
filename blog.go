package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/russross/blackfriday"
)

func main() {
	var (
		mux   = pat.New()
		port  = flag.String("p", "3000", "address to bind the server on")
		dev   = flag.Bool("dev", false, "enable/disable dev mode")
		posts = allPosts()
		t, _  = template.ParseGlob("./templates/*.tmpl")
	)
	flag.Parse()

	if *dev {
		mux.Get("/new_post", func(w http.ResponseWriter, r *http.Request) {
			t.ExecuteTemplate(w, "new.tmpl", nil)
		})
		mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			path := "./public/img/" + header.Filename
			img, err := os.Create(path)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer img.Close()

			io.Copy(img, file)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(path))
		})
		mux.Post("/new_post", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.FormValue("title"))
			p, err := CreatePost(r.FormValue("title"), r.FormValue("body"))
			if err != nil {
				http.Error(w, "Error making post.", 500)
				return
			}
			posts = allPosts()
			t.ExecuteTemplate(w, "show.tmpl", p)
		})
		mux.Get("/markdown_preview", func(w http.ResponseWriter, r *http.Request) {
			ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
			if _, ok := err.(websocket.HandshakeError); ok {
				http.Error(w, "Not a websocket handshake", 400)
				return
			} else if err != nil {
				return
			}

			for {
				messageType, message, err := ws.ReadMessage()
				if err != nil {
					return
				}
				if err := ws.WriteMessage(messageType, blackfriday.MarkdownCommon(message)); err != nil {
					return
				}
			}
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
		t.ExecuteTemplate(w, "show.tmpl", p)
	})

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.tmpl", posts)
	})

	log.Printf("Starting server on %s\n", *port)
	handler := handlers.LoggingHandler(os.Stdout, mux)
	handler = handlers.CompressHandler(handler)
	log.Fatal(http.ListenAndServe(":"+*port, handler))
}
