package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("F:\\Client\\templates\\startbootstrap-simple-sidebar-gh-pages1\\index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func main() {
	// fileServer := http.FileServer(http.Dir("./index"))
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8545", nil)
	// http.HandleFunc("/form", formHandler)
	// http.HandleFunc("/hello", helloHandler)

	// fmt.Printf("Starting server at port 8080\n")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal(err)
	// }
}
