/*
.
.
SERVER CODE
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
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/serverdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("succussfuly connect")
	///////////////////////////////////////////////////network programing
	fmt.Println("server listening on 3000")
	listener, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	// listening for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		buffer := make([]byte, 1400)
		dataSize, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("connection closed")
			return
		}
		// the actual message
		data := buffer[:dataSize]
		fmt.Println("received message: ", string(data))
		// listen to connections in another gorutine
		go listenConnection(conn, db)
	}
}

// listening for messages from connection
func listenConnection(conn net.Conn, db *sql.DB) {
	for {
		choose, _ := read_fromClient(conn)
		if choose == "1" {
			_Title, _ := read_fromClient(conn)
			albID, err := addAlbum(db, Album{
				Title: _Title,
			})
			fmt.Printf("ID of added album: %v \n", albID)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("succussfuly insert  sir ")
		}
	}
}
func read_fromClient(conn net.Conn) (string, error) {

	buffer_command := make([]byte, 1400)
	comm_Size, err := conn.Read(buffer_command)
	if err != nil {
		fmt.Println("connection closed")
		return "err", nil
	}
	command := string(buffer_command[:comm_Size])
	println("read " + command)
	return command, nil
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
