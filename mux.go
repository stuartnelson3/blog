package main

import(
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "text/template"
    "path/filepath"
    "io/ioutil"
)

type Post struct {
    Title       string `json:"title"`
    Body        string `json:"body"`
    CreatedAt   string `json:"createdAt"`
    Slug        string `json:"slug"`
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", rootHandler)
    r.HandleFunc("/{post}", postHandler)
    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    absPath, _ := filepath.Abs("on-formatting.json")
    data, _ := ioutil.ReadFile(absPath)
    out := &Post{}
    json.Unmarshal(data, out)
    t, _ := template.New("root").Parse(`{{define "T"}}{{.}}{{end}}`)
    _ = t.ExecuteTemplate(w, "T", out.Body)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    post := vars["post"]
    fmt.Fprintf(w, "Hello, %q", post)
}
