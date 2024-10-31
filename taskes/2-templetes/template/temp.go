package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}
type TodoPage struct {
	PageTitle string
	Todos     []Todo
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "welcom to index  page ")
}
func todosHandler(w http.ResponseWriter, r *http.Request) {

	data := TodoPage{
		PageTitle: "all todoes",
		Todos: []Todo{
			{Title: "learn go", Done: false},
			{Title: "learn python", Done: true},
			{Title: "learn javascript", Done: true},
		},
	}
	t, err := template.ParseFiles("./template/index.html")

	if err != nil {
		log.Fatal("error in counted  while pasring the templete :", err)
	}
	t.Execute(w, data)
}
func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/todos", todosHandler)
	http.ListenAndServe(":8080", nil)

}
