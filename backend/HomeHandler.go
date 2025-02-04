package backend

import (
	"html/template"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}
	apiArtist := "https://groupietrackers.herokuapp.com/api/artists"

	artists, err := fetchArtists(apiArtist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}
	tmpl.Execute(w, artists)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, er := template.ParseFiles("templates/404.html")
	w.WriteHeader(http.StatusNotFound)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}
	tmpl.Execute(w, nil)
}
