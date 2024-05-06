package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const PORT = ":8080"

var tpl *template.Template

func init() {
	// Parsing all html files
	tpl = template.Must(template.New("").ParseGlob("views/**/*.html"))
	// Listen statics files
	fs := http.FileServer(http.Dir("public/assets"))
	http.Handle("/public/", http.StripPrefix("/public/assets", fs))
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func main() {

	// localHost link
	fmt.Println("serving at : http://localhost" + PORT)

	http.HandleFunc("/", Home)
	// Start forum
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalln(err)
	}
}
