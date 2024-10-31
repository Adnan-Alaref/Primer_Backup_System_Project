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

	_ "github.com/go-sql-driver/mysql"
)

type Album struct {
	ID    int64
	Title string
}

func main() {
	///////////////////////////////////////////////////connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/clientdb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	fmt.Println("succussfuly connect")
	///////////////////////////////////////////////////network programing
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("Hello Server!"))
	if err != nil {
		log.Fatalln(err)
	}
	for {
		var choose int
		fmt.Println("enter 1 for insert & 2  to select by artist Name")
		fmt.Scan(&choose)
		if choose == 1 {
			fmt.Println("enter the Title ")
			var _Title string
			fmt.Scan(&_Title)
			albID, err := addAlbum(db, Album{
				Title: _Title,
			})
			fmt.Printf("ID of added album: %v\n", albID)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("succussfuly insert  sir ")
			fmt.Println("send command and data  to servser ")
			_, err = conn.Write([]byte("1"))
			_, err = conn.Write([]byte(_Title))
			if err != nil {
				log.Fatalln(err)
			}
		}

		////////////////////////////////////////send data to server
		// buffer := make([]byte, 1400)
		// dataSize, err := conn.Read(buffer)
		// if err != nil {
		//     fmt.Println("connection closed")
		//     return
		// }
		// data := buffer[:dataSize]
		// fmt.Println("received message: ", string(data))
	}
	fmt.Println("succussfuly end program ")

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
