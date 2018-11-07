package main

import (
	"database/sql"
	"net/http"

	"html/template"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Story    string
	DBStatus bool
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := Page{Story: "Something intense page"}
		if story := r.FormValue("story"); story != "" {
			p.Story = story
		}
		p.DBStatus = db.Ping() == nil

		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":8080", nil)
}
