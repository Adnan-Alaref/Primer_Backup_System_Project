package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	ID    int64
	err   error
	db    *sql.DB
	templ *template.Template
)

type PersonData struct {
	Id          int
	Username    string
	Password    string
	Description string
	Created_at  mysql.NullTime
	Updated_at  mysql.NullTime
	AdminD      int
}

var templates *template.Template

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/crud_db")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/create", createUser).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8090", nil)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {

	templates.ExecuteTemplate(w, "index.html", nil)
}
func createUser(res http.ResponseWriter, req *http.Request) {
	//insert into db
	stmt, err := db.Prepare("INSERT Crud_tb SET username=?, password=?, description=?, created_at=?, updated_at=?, is_admin=?")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	if req.Method != "POST" {
		templ.ExecuteTemplate(res, "create.html", nil)
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	describe := req.FormValue("description")
	isAdmin := req.FormValue("admin")
	createdAt := time.Now()
	updatedAt := time.Now()
	var admin_chk int
	if isAdmin == "on" {
		admin_chk = 1
	} else {
		admin_chk = 0
	}
	var user string
	err = db.QueryRow("SELECT username FROM Crud_tb WHERE username=?", username).Scan(&user)
	switch {
	//username is available
	case err == sql.ErrNoRows:
		securedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		rs, err := stmt.Exec(username, securedPassword, describe, createdAt, updatedAt, admin_chk)
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		id, err := rs.LastInsertId()
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		ID = getID(id).(int64)
		res.Write([]byte("user successfully created!"))
		fmt.Println("user: ", username, " with ID: ", id, " successfully created!")
		return
	case err != nil:
		http.Error(res, err.Error(), 500)
		return
	default:
		http.Redirect(res, req, "/create", 301)
	}
}
func readUser(res http.ResponseWriter, req *http.Request) {
	// query
	rows, err := db.Query("SELECT * FROM Crud_tb")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	var id int
	var username string
	var password string
	var describe string
	var created_at mysql.NullTime
	var updated_at mysql.NullTime
	var isAdmin int
	/*if req.Method != "POST" {
	}*/
	var ps []PersonData
	for rows.Next() {
		err = rows.Scan(&id, &username, &password, &describe, &created_at, &updated_at, &isAdmin)
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		ps = append(ps, PersonData{Id: id, Username: username, Password: password, Description: describe, Created_at: created_at, Updated_at: updated_at, AdminD: isAdmin})
		//return
	}
	templ.ExecuteTemplate(res, "read.html", ps)
}
func getID(id int64) interface{} {
	return id
}
