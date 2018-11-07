package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"html/template"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Name     string
	DBStatus bool
}

type SearchResult struct {
	Title  string
	Author string
	Year   string
	ID     string
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := Page{Name: "Vello"}
		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}
		p.DBStatus = db.Ping() == nil

		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//db.Close()
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results := []SearchResult{
			SearchResult{"Interesting Title", "Curious Author", "1999", "12341234"},
			SearchResult{"Boring Title", "Dull Author", "2010", "12322312"},
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
