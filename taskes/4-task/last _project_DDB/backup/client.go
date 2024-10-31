/*
.
.
CLIENT CODE
.
.
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Album struct {
	ID    int64
	Title string
	// Artist string
	// Price string
}

var db *sql.DB
var err error

func Connect_DB() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/recordings")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("succussfuly connect")
}
func main() {
	///////////////////////////////////////////////////connection
	Connect_DB()

	///////////////////////////////////////////////////network programing
	// conn, err := net.Dial("tcp", "192.168.43.103:3000") 
	listener, err := net.Listen("tcp","0.0.0.0:8000")
	if err!=nil{
		log.Fatal(err)
	} 
	conn,err := listener.Accept()
	if err!=nil{
		log.Fatal(err)
	}
    Say_hallo := make([]byte, 1400)
    datasize5, _:= conn.Read(Say_hallo)
	massgae := string(Say_hallo[:datasize5])
	fmt.Println(massgae)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// _, err = conn.Write([]byte("Hello Server!"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	for {
		choose, err := ReadFromServer(conn)
		if err != nil {
			fmt.Println("Connection closed..!")
			return
		} else if choose == "1" {
			_Title, _ := ReadFromServer(conn)
			albID, err := addAlbum(db, Album{
				Title: _Title,
			})
			fmt.Printf("ID of added album: %v \n", albID)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("successfully insert  sir ")
		} else if choose == "2" {
			_id, _ := ReadFromServer(conn)
			id, _ := strconv.Atoi(_id)
			DeleteAlbum(id)
		} else if choose == "3" {
			_id, _ := ReadFromServer(conn)
			_Title, _ := ReadFromServer(conn)
			id, _ := strconv.Atoi(_id)
			Update_Album(id, _Title)
		}
	}
}

func ReadFromServer(conn net.Conn) (string, error) {
	buffer_command := make([]byte, 1400)
	dataSize, err := conn.Read(buffer_command)
	if err != nil {
		log.Fatal(err)
	}
	command := string(buffer_command[:dataSize])
	return command, err
}

func addAlbum(db *sql.DB, alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title) VALUES (?)", alb.Title)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
// func GetAllAlbums() []Album {
// 	row, err := db.Query("select * from album")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	alb := Album{}      //Create instanc from Album
// 	albums := []Album{} //Create Array of Album
// 	for row.Next() {
// 		err := row.Scan(&alb.ID, &alb.Title)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		albums = append(albums, alb)
// 	}
// 	return albums
// }

func Update_Album(id int, title string) {
	statment, err := db.Prepare("update recordings.album set Title=? where ID=?")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statment.Exec(title, id)
	if err != nil {
		log.Fatal(err)
	}
	affectedRow, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Query OK, %d rows affected.\n", affectedRow)
		fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")

	}
}
func DeleteAlbum(id int) {
	statment, err := db.Prepare("delete from album where ID=?")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statment.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	affectedRow, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Query OK, %d rows affected.\n", affectedRow)
		fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")
	}
}
