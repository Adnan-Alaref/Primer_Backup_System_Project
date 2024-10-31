package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {

	///////////////////connection section
	fmt.Println("Go MySQL Tutorial")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/recordings")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	fmt.Println("succussfuly connect")
	//////////////////////////////////////////////////my  work  sir (:

	for i := 0; i < 3; i++ {
		var choose int
		fmt.Println("enter 1 for insert & 2  to select by artist Name")
		fmt.Println("enter 3 for select by id & 4  for update price")
		fmt.Println("enter 5 for delete  by id \n")

		fmt.Scan(&choose)
		if choose == 1 {
			fmt.Println("enter the Title & Artist & Price")
			var _Title, _Artist string
			var _Price float32
			fmt.Scan(&_Title, &_Artist, &_Price)
			albID, err := addAlbum(db, Album{
				Title:  _Title,
				Artist: _Artist,
				Price:  _Price,
			})
			fmt.Printf("ID of added album: %v\n", albID)
			if err != nil {
				panic(err.Error())
			}
			// be careful deferring Queries if you are using transactions

			fmt.Println("succussfuly insert  sir ")
		}
		/////////////////////////////////////select section
		if choose == 2 {
			fmt.Println("enter the name  of Artist and  i shaw all related  row ")
			var ArtistNmae1, ArtistNmae2 string
			fmt.Scan(&ArtistNmae1, &ArtistNmae2)

			Albumbs, _ := albumsByArtist(db, ArtistNmae1+" "+ArtistNmae2)
			for _, val := range Albumbs {
				fmt.Printf("Id is %v and Title  is %v and Artist  %v \n\n", val.ID, val.Title, val.Artist)
			}
			fmt.Println("succussfuly select by artist sir ")
		}
		if choose == 3 {
			fmt.Println("enter the id  and  i shaw all related  row ")
			var Id int64
			fmt.Scan(&Id)
			alb, err := albumByID(db, Id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Album found: %v\n", alb)

			fmt.Println("succussfuly select  by id sir ")
		}
		if choose == 4 {

			updatePrice(db)

			fmt.Println("succussfuly update Price  sir ")
		}
		if choose == 5 {

			Deletealbum(db)

			fmt.Println("succussfuly delete  sir ):  ")
		}

	}
	////////////////////////////////////////end
	fmt.Println("succussfuly end program ")
}
func albumsByArtist(db *sql.DB, name string) ([]Album, error) {

	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(db *sql.DB, id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(db *sql.DB, alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func updatePrice(db *sql.DB) {
	stmt, e := db.Prepare("update album set price=? where id=?")
	checkErr(e)
	fmt.Println("enter new price  and id want change ")
	var newPrice float32
	var id int
	fmt.Scan(&newPrice, &id)
	// execute
	res, e := stmt.Exec(newPrice, id)
	checkErr(e)

	a, e := res.RowsAffected()
	checkErr(e)

	fmt.Println(a) // 1
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
func Deletealbum(db *sql.DB) {

	stmt, e := db.Prepare("DELETE FROM album WHERE id=?")
	checkErr(e)
	fmt.Println("enter  id  wnt delete it ")
	var id int
	fmt.Scan(&id)
	// execute
	res, e := stmt.Exec(id)
	checkErr(e)

	a, e := res.RowsAffected()
	checkErr(e)

	fmt.Println(a)
}
