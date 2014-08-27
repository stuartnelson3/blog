package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/russross/blackfriday"
)

var (
	posts = allPosts()
	t, _  = template.ParseGlob("./templates/*.tmpl")
)

func main() {
	var (
		mux  = pat.New()
		port = flag.String("p", "3000", "address to bind the server on")
		dev  = flag.Bool("dev", false, "enable/disable dev mode")
	)
	flag.Parse()

	if *dev {
		mux.Get("/new_post", func(w http.ResponseWriter, r *http.Request) {
			t.ExecuteTemplate(w, "new.tmpl", nil)
		})
		mux.Post("/upload", errorHandler(upload))
		mux.Post("/new_post", errorHandler(create))
		mux.Get("/markdown_preview", errorHandler(markdownPreview))
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

func errorHandler(fn func(http.ResponseWriter, *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if statusCode, err := fn(w, r); err != nil {
			http.Error(w, err.Error(), statusCode)
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) (int, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return 500, err
	}
	defer file.Close()

	imgPath, err := CreateImage(file, header)
	if err != nil {
		return 500, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(imgPath))
	return 200, nil
}

func markdownPreview(w http.ResponseWriter, r *http.Request) (int, error) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return 500, errors.New("Not a websocket handshake")
	} else if err != nil {
		return 500, err
	}

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			return 500, err
		}
		if err := ws.WriteMessage(messageType, blackfriday.MarkdownCommon(message)); err != nil {
			return 500, err
		}
	}
}

func create(w http.ResponseWriter, r *http.Request) (int, error) {
	p, err := CreatePost(r.FormValue("title"), r.FormValue("body"))
	if err != nil {
		return 500, errors.New("Error making post.")
	}
	posts = allPosts()
	t.ExecuteTemplate(w, "show.tmpl", p)
	return 200, nil
}
