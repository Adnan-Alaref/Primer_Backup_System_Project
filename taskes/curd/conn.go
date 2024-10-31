package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
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

func main() {

	///////////////////connection section
	fmt.Println("Go MySQL Tutorial")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/crud_db")
	checkErr(err)
	templ, err = templ.ParseGlob("templates/*.html")
	checkErr(err)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	fmt.Println("succussfuly connect")
	//////////////////////////////////////////////////my  work  sir (:
	http.HandleFunc("/", index)
	http.HandleFunc("/create", createUser)
	http.HandleFunc("/read", readUser)
	http.HandleFunc("/update", updateUser)
	http.HandleFunc("/delete", deleteUser)
	http.Handle("/assets/", http.FileServer(http.Dir("."))) //serve other files in assets dir
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fmt.Println("server running on port :8080")
	http.ListenAndServe(":8080", nil)

	////////////////////////////////////////end
	println("succussfuly end program ")
}
func index(res http.ResponseWriter, req *http.Request) {
	templ.ExecuteTemplate(res, "index.html", nil)
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
func updateUser(res http.ResponseWriter, req *http.Request) {
	//select id's
	rows, err := db.Query("SELECT id FROM Crud_tb")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	var user = req.FormValue("ids")
	var newUsername = req.FormValue("username")
	var ps []PersonData
	id, err := strconv.Atoi(user)
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		ps = append(ps, PersonData{Id: id})
	}
	stmt, err := db.Prepare("UPDATE Crud_tb SET username=?, updated_at=? WHERE id=?")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	rs, err := stmt.Exec(newUsername, time.Now(), id)
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	affect, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	if req.Method != "POST" {
		templ.ExecuteTemplate(res, "update.html", ps)
		return
	}
	fmt.Println("row :", affect, " affected")
}
func deleteUser(res http.ResponseWriter, req *http.Request) {
	//select id's
	rows, err := db.Query("SELECT id FROM Crud_tb")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusNoContent)
		return
	}
	var user = req.FormValue("ids")
	var ps []PersonData
	id, err := strconv.Atoi(user)
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			http.Error(res, "there was an error", http.StatusInternalServerError)
			return
		}
		ps = append(ps, PersonData{Id: id})
	}
	// delete
	stmt, err := db.Prepare("delete from Crud_tb where id=?")
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	rs, err := stmt.Exec(id)
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	affect, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(res, "there was an error", http.StatusInternalServerError)
		return
	}
	if req.Method != "POST" {
		templ.ExecuteTemplate(res, "delete.html", ps)
		return
	}
	fmt.Println("row :", affect, " affected")
}
func getID(id int64) interface{} {
	return id
}
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
