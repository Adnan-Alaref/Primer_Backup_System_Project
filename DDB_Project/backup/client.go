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
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Album struct {
	ID        int64
	Title     string
	Artist_id int
	Price     int
}

type artist struct {
	ID    int64
	name  string
	email string
}

var db *sql.DB
var err error

// conect to the database
func Connect_DB() {
	db, err = sql.Open("mysql", "root:omer8520@tcp(127.0.0.1:3306)/recordings")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("succussfuly connect")
}

func main() {
	///////////////////////////////////////////////////connection database
	Connect_DB()

	///////////////////////////////////////////////////network programing
	// conn, err := net.Dial("tcp", "192.168.43.103:3000")
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println("Successfuly conect to the Master :)")
	fmt.Println()
	Say_hallo := make([]byte, 1400)
	datasize5, _ := conn.Read(Say_hallo)
	massgae := string(Say_hallo[:datasize5])
	fmt.Println(massgae)
	_, err = conn.Write([]byte("Hello Server! The Puckaup Is Ready :)"))
	if err != nil {
		log.Fatalln(err)
	}

	////////////////  insert delete and update looop  \\\\\\\\\\\\\\\\
	for {
		choose, err := ReadFromServer(conn)
		if err != nil {
			fmt.Println("Connection closed..!", err.Error())
			return
		} else if choose == "1" {

			fmt.Println("Start inserting")
			inserteddatat, _ := ReadFromServer(conn)
			data := strings.Split(inserteddatat, ",")
			if data[0] == "album" {

				_Title := data[1]
				_artist := data[2]
				_price := data[3]
				intprice, _ := strconv.Atoi(_price)
				artistid, _ := strconv.Atoi(_artist)
				albID, err := addAlbum(db, Album{
					Title:     _Title,
					Artist_id: artistid,
					Price:     intprice,
				})
				fmt.Printf("ID of added album: %v \n", albID)
				if err != nil {
					panic(err.Error())
				}
				fmt.Println("successfully insert  sir ")
			} else if data[0] == "artist" {
				_name := data[1]
				_email := data[2]
				artID, err := addArtist(db, artist{
					name:  _name,
					email: _email,
				})
				fmt.Printf("ID of added artist: %v \n", artID)
				if err != nil {
					panic(err.Error())
				}
				fmt.Println("successfully insert  sir ")
			}
		} else if choose == "2" {
			fmt.Println("Start deleting")
			inserteddatat, _ := ReadFromServer(conn)
			data := strings.Split(inserteddatat, ",")
			if data[0] == "album" {
				// _id, _ := ReadFromServer(conn)
				id, _ := strconv.Atoi(data[1])
				DeleteAlbum(id)
			} else if data[0] == "artist" {
				id, _ := strconv.Atoi(data[1])
				DeleteArtist(id)
			}
			fmt.Println("successfully Deleted  sir ")
		} else if choose == "3" {
			fmt.Println("Start Updating")
			inserteddatat, _ := ReadFromServer(conn)
			data := strings.Split(inserteddatat, ",")
			if data[0] == "album" {
				// updataeddatat, _ := ReadFromServer(conn)
				// data := strings.Split(updataeddatat, ",")
				_id := data[1]
				_Title := data[2]
				_artist := data[3]
				_price := data[4]
				newprice, _ := strconv.Atoi(_price)
				id, _ := strconv.Atoi(_id)
				Update_Album(id, _Title, _artist, newprice)
			} else if data[0] == "artist" {
				// updataeddatat, _ := ReadFromServer(conn)
				// data := strings.Split(updataeddatat, ",")
				_id := data[1]
				_name := data[2]
				_email := data[3]
				id, _ := strconv.Atoi(_id)
				Update_artist(id, _name, _email)
			}
			fmt.Println("successfully Updated sir ")
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

//Inserting funcctions
func addAlbum(db *sql.DB, alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title,artist_id,price) VALUES (?,?,?)", alb.Title, alb.Artist_id, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
func addArtist(db *sql.DB, art artist) (int64, error) {
	result, err := db.Exec("INSERT INTO artist (name,email) VALUES (?,?)", art.name, art.email)
	if err != nil {
		return 0, fmt.Errorf("addArtist: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addArtist: %v", err)
	}
	return id, nil
}

//Updated Functanis
func Update_Album(id int, title string, artist_id string, price int) {
	statment, err := db.Prepare("update recordings.album set title=?, artist_id=?, price=? where id=?")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statment.Exec(title, artist_id, price, id)
	if err != nil {
		log.Fatal(err)
	}
	affectedRow, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("The statment afected %d rows.\n", affectedRow)
		// fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")

	}
}
func Update_artist(id int, name string, email string) {
	statment, err := db.Prepare("update recordings.artist set name=?, email=? where id=?")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statment.Exec(name, email, id)
	if err != nil {
		log.Fatal(err)
	}
	affectedRow, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("The statment afected %d rows.\n", affectedRow)
		// fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")

	}
}

// Deletign finctions
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
		fmt.Printf("The statment afected %d rows.\n", affectedRow)
		// fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")
	}
}
func DeleteArtist(id int) {
	statment1, err := db.Prepare("delete from album where artist_id=?")
	if err != nil {
		log.Fatal(err)
	}
	r1, err := statment1.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	affectedRow1, err := r1.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("The statment afected %d rows.\n", affectedRow1)
		// fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")
	}
	statment, err := db.Prepare("delete from artist where ID=?")
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
		fmt.Printf("The statment afected %d rows.\n", affectedRow)
		// fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")
	}

}
