/*
.
.
 Author : Adnan Alarf
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
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Album struct {
	ID       int64
	Title    string
	ArtistID int
	Price    int64
}

type Artist struct {
	ID    int64
	Name  string
	Email string
}

var conn_backup net.Conn

func main() {
	/*
		.Connect to database
	*/
	dp_connection()
	/*
		.Connect to backup
	*/
	conn_backup, _ = net.Dial("tcp", "192.168.87.181:8000")
	message := make([]byte, 1400)
	_, _ = conn_backup.Write([]byte("Hello Backup,Are You There!"))
	datasize, _ := conn_backup.Read(message)
	str := string(message[:datasize])
	fmt.Println(str)

	fmt.Println("server listening on 3000")
	/*
		. Accept any connection
	*/
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

func dp_connection() {
	dpDriver := "mysql"
	dpUser := "root"
	dpPass := ""
	dpName := "recordings"

	db, err = sql.Open(dpDriver, dpUser+":"+dpPass+"@tcp(127.0.0.1:3306)/"+dpName)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	fmt.Println("succussfuly connect")
}

//Drop Played , Track ,Album ,Artist
// Listening for messages from connection
func listenConnection(conn net.Conn, db *sql.DB) {
	for {
		/*
			. receive type of command in choose varible
		*/
		choose, err := read_fromClient(conn)
		//send data for backup
		write_to_backup(conn_backup, choose)
		if err != nil {
			fmt.Println("Connection closed..!", err.Error())
			return
		} else if choose == "1" {
			/*
			  Insert
			*/
			/*
				. create slice from album
			*/
			var attributes []string = make([]string, 4)
			/*
				. All data from client
			*/
			data, _ := read_fromClient(conn)
			/*.make slice to receve data and split it
			 */
			attributes = strings.Split(data, ",")
			//tack table name
			table_name := attributes[0]
			/*
				.make conditions over table name
				.there are two tables album & artist
			*/
			if table_name == "album" {
				//convert string to int
				artist_id, _ := strconv.Atoi(attributes[2])
				price, _ := strconv.Atoi(attributes[3])
				InsertData(Album{
					Title:    attributes[1],
					ArtistID: artist_id,
					Price:    int64(price),
				})
				if err != nil {
					panic(err.Error())
				}
				/*
				   . Send data to backup
				*/
				write_to_backup(conn_backup, data)
			} else if table_name == "artist" {
				InsertData(Artist{
					Name:  attributes[1],
					Email: attributes[2],
				})
				if err != nil {
					panic(err.Error())
				}
				/*
					. Send data to backup
				*/
				write_to_backup(conn_backup, data)
			}
		} else if choose == "2" {
			/*
				Delete
			*/
			tab_name, _ := read_fromClient(conn)
			if tab_name == "album" {
				arr := []Album{}
				arr = GetAllAlbums()
				var Allids string
				_, _ = conn.Write([]byte(strconv.Itoa(len(arr))))
				for i := 0; i < len(arr); i++ {
					id := arr[i].ID
					subid := strconv.Itoa(int(id))
					res := subid + ";"
					if i == len(arr)-1 {
						res = res[:len(res)-1]
						Allids += res
					} else {
						Allids += res
					}
				}
				_, _ = conn.Write([]byte(Allids))
				// fmt.Println("From Function", arr)
				// fmt.Println("From All rows ", Allids)

				//receive id from clients
				_id, _ := read_fromClient(conn)
				id, _ := strconv.Atoi(_id)
				Delete_Data(int64(id), Album{})
				/*
					. Send data to backup
				*/
				data := tab_name + "," + _id
				write_to_backup(conn_backup, data)
			} else if tab_name == "artist" {
				arr1 := []Artist{}
				arr1 = GetAllArtists()
				var Allids string

				_, _ = conn.Write([]byte(strconv.Itoa(len(arr1))))
				for i := 0; i < len(arr1); i++ {
					id := arr1[i].ID
					subid := strconv.Itoa(int(id))
					res := subid + ";"
					if i == len(arr1)-1 {
						res = res[:len(res)-1]
						Allids += res
					} else {
						Allids += res
					}
				}
				_, _ = conn.Write([]byte(Allids))
				// fmt.Println("From Function", arr1)
				// fmt.Println("From All rows ", Allids)

				//receive id from clients
				_id, _ := read_fromClient(conn)
				id, _ := strconv.Atoi(_id)
				Delete_Data(int64(id), Artist{})
				/*
					. Send data to backup
				*/
				data := tab_name + "," + _id
				write_to_backup(conn_backup, data)
			}
		} else if choose == "3" {
			/*
			 Updatefgh
			*/
			/*
				Recieve New Data From Client
			*/
			var newdata []string = make([]string, 5)
			// read from cliect
			new_data, _ := read_fromClient(conn)
			//make slice
			newdata = strings.Split(new_data, ",")
			table_name := newdata[0]
			//check over tablename
			if table_name == "album" {
				id, _ := strconv.Atoi(newdata[1])
				artist_id, _ := strconv.Atoi(newdata[3])
				price, _ := strconv.Atoi(newdata[4])
				Update_Data(Album{
					ID:       int64(id),
					Title:    newdata[2],
					ArtistID: artist_id,
					Price:    int64(price),
				})
				/*
					. Send data to backup
				*/
				write_to_backup(conn_backup, new_data)
			} else if table_name == "artist" {
				id, _ := strconv.Atoi(newdata[1])
				Update_Data(Artist{
					ID:    int64(id),
					Name:  newdata[2],
					Email: newdata[3],
				})
				/*
					. Send data to backup
				*/
				write_to_backup(conn_backup, new_data)
			}
		} else if choose == "4" {
			/*
			  Select
			*/
			table_name, _ := read_fromClient(conn)
			if table_name == "album" {
				arr := []Album{}
				arr = GetAllAlbums()
				// fmt.Println(len(arr))
				write_fromClient(strconv.Itoa(len(arr)), conn)
				var Allrows string
				for i := 0; i < len(arr); i++ {
					id := arr[i].ID
					title := arr[i].Title
					artistid := arr[i].ArtistID
					price := arr[i].Price
					subslic := []string{strconv.Itoa(int(id)), title, strconv.Itoa(int(artistid)), strconv.Itoa(int(price))}
					result := strings.Join(subslic, ",")
					row := result + ";"
					if i == len(arr)-1 {
						row = row[:len(row)-1]
						Allrows += row
					} else {
						Allrows += row
					}
				}
				_, _ = conn.Write([]byte(Allrows))
				fmt.Printf("Query OK, %d rows in set (0.00 sec).\n", len(arr))
				fmt.Println("Data In Table ", GetAllAlbums())
				// fmt.Println("From All rows ", Allrows)
			} else if table_name == "artist" {
				arr1 := []Artist{}
				arr1 = GetAllArtists()
				// fmt.Println(len(arr1))
				write_fromClient(strconv.Itoa(len(arr1)), conn)
				var Allrows string
				for i := 0; i < len(arr1); i++ {
					id := arr1[i].ID
					name := arr1[i].Name
					email := arr1[i].Email
					subslic := []string{strconv.Itoa(int(id)), name, email}
					result := strings.Join(subslic, ",")
					row := result + ";"
					if i == len(arr1)-1 {
						row = row[:len(row)-1]
						Allrows += row
					} else {
						Allrows += row
					}
				}
				_, _ = conn.Write([]byte(Allrows))
				fmt.Printf("Query OK, %d rows in set (0.00 sec).\n", len(arr1))
				fmt.Println("Data In Table ", GetAllArtists())
				// fmt.Println("From All rows ", Allrows)
			}
		}
	}
}
func read_fromClient(conn net.Conn) (string, error) {
	buffer_command := make([]byte, 1400)
	comm_Size, err := conn.Read(buffer_command)
	if err != nil {
		// fmt.Println("Connection closed")
		return "", err
	}
	command := string(buffer_command[:comm_Size])
	fmt.Println("read " + command)
	return command, nil
}

func write_to_backup(conn net.Conn, str string) error {
	_, err := conn.Write([]byte(str))
	if err != nil {
		//fmt.Println("Connection closed")
		return err
	}
	return nil
}
func write_fromClient(data string, conn net.Conn) {
	conn.Write([]byte(data))
}

func InsertData(data interface{}) {
	if alb, ok := data.(Album); ok {
		result, err1 := db.Exec("INSERT INTO recordings.album (title,artist_id,price) VALUES (?,?,?)", alb.Title, alb.ArtistID, alb.Price)
		if err1 != nil {
			log.Fatal(err1.Error())
		}
		alID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err.Error())
		}
		affectedRow, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID Of Added Album: %v \n", alID)
		fmt.Printf("Query OK, %d rows affected.\n", affectedRow)
	} else if art, ok := data.(Artist); ok {
		result, err1 := db.Exec("INSERT INTO recordings.artist (name,email) VALUES (?,?)", art.Name, art.Email)
		if err1 != nil {
			log.Fatal(err1.Error())
		}
		artID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err.Error())
		}
		affectedRow1, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of added artist: %v \n", artID)
		fmt.Printf("Query OK, %d rows affected.\n", affectedRow1)
	}
}

func GetAllAlbums() []Album {
	row, err := db.Query("select * from album")
	if err != nil {
		log.Fatal(err)
	}
	alb := Album{}      //Create instanc from Album
	albums := []Album{} //Create Array of Album
	for row.Next() {
		err := row.Scan(&alb.ID, &alb.Title, &alb.ArtistID, &alb.Price)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, alb)
	}
	return albums
}

func GetAllArtists() []Artist {
	row, err := db.Query("select * from artist")
	if err != nil {
		log.Fatal(err)
	}
	art := Artist{}       //Create instanc from Artist
	artists := []Artist{} //Create Array of Artist
	for row.Next() {
		err := row.Scan(&art.ID, &art.Name, &art.Email)
		if err != nil {
			log.Fatal(err)
		}
		artists = append(artists, art)
	}
	return artists
}
func Update_Data(data interface{}) {
	if alb, ok := data.(Album); ok {
		statment, err := db.Prepare("update recordings.album set title=?, artist_id=? ,price=? where ID=?")
		if err != nil {
			log.Fatal(err)
		}
		r, err := statment.Exec(alb.Title, alb.ArtistID, alb.Price, alb.ID)
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
	} else if art, ok := data.(Artist); ok {
		statment, err := db.Prepare("update recordings.artist set name =?,email=? where ID=?")
		if err != nil {
			log.Fatal(err)
		}
		r, err := statment.Exec(art.Name, art.Email, art.ID)
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
}

func Delete_Data(id int64, data interface{}) {

	if _, ok := data.(Album); ok {
		/*
		 .Delete From album table
		*/
		statment, err := db.Prepare("delete from recordings.album where ID=?")
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
	} else if _, ok := data.(Artist); ok {
		/*
		 .Delete From artist table
		*/
		/*
			. First
			. Start Delete From Album
		*/
		statment1, err := db.Prepare("delete from recordings.album where artist_id=?")
		if err != nil {
			log.Fatal(err.Error())
		}
		s, _ := statment1.Exec(id)
		_, _ = s.RowsAffected()
		fmt.Printf("Query OK, %d album rows affected.\n", 1)
		/*
			.Second
			. Start Delete From Artist
		*/
		statment, err := db.Prepare("delete from recordings.artist where ID=?")
		if err != nil {
			log.Fatal(err)
		}
		r, err1 := statment.Exec(id)
		if err1 != nil {
			log.Fatal(err1)
		}
		affectedRow, err2 := r.RowsAffected()
		if err2 != nil {
			log.Fatal(err2)
		} else {
			fmt.Printf("Query OK, %d artist rows affected.\n", affectedRow)
			fmt.Println("Rows matched: 1  Changed: 1  Warnings: 0")
		}
	}
}
