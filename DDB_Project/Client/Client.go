/*
.
.
CLIENT CODE
.
.
*/
package main

import (
	// "fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	artist "example.com/module/Artist"
	"example.com/module/tables"
	"github.com/gorilla/mux"
)

var err error
var conn net.Conn

func render(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, data)
}

type IDS struct {
	IDs []int
}
type Albums struct {
	AllAlbums []tables.Album
}
type Artists struct {
	AllArtists []artist.Artist
}

func main() {
	r := mux.NewRouter()
	//ConnectToMaster
	conn, err = net.Dial("tcp", "192.168.43.11:3000")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("Hello Server!"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		r.HandleFunc("/", homehandler)

		r.HandleFunc("/create_album", InsertAlbum).Methods("GET", "POST")
		r.HandleFunc("/create_artist", InsertArtist).Methods("GET", "POST")

		r.HandleFunc("/delete_album", DeleteAlbum).Methods("GET", "POST")
		r.HandleFunc("/delete_artist", DeleteArtist).Methods("GET", "POST")

		r.HandleFunc("/update_album", UpdateAlbum).Methods("GET", "POST")
		r.HandleFunc("/update_artist", UpdateArtist).Methods("GET", "POST")

		r.HandleFunc("/select_album", SelectAlbum).Methods("GET", "POST")
		r.HandleFunc("/select_artist", SelectArtist).Methods("GET", "POST")

		http.ListenAndServe(":8080", r)
	}

}
func homehandler(w http.ResponseWriter, req *http.Request) {
	render(w, "F:\\Client\\templates\\index.html", nil)
}
func InsertAlbum(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {

		render(w, "F:\\Client\\templates\\insert_album.html", nil)
	}
	if req.Method == "POST" {

		price, _ := strconv.Atoi(req.PostFormValue("price"))
		artist_id, _ := strconv.Atoi(req.PostFormValue("artist_id"))
		msg_album := &tables.Album{

			TableName: "album",
			ID:        1, //
			Artist_ID: int64(artist_id),
			Title:     req.PostFormValue("title"),
			Price:     int64(price),
		}
		if msg_album.Validate() == false {
			render(w, "F:\\Client\\templates\\insert_album.html", msg_album)
			return
		}
		_, err = conn.Write([]byte("1"))
		if err != nil {
			log.Fatal(err.Error())
		}

		str_price := strconv.Itoa(price)
		str_artist_id := strconv.Itoa(artist_id)
		myslice := []string{msg_album.TableName, msg_album.Title, str_artist_id, str_price}
		data := strings.Join(myslice, ",")
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
func InsertArtist(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {

		render(w, "F:\\Client\\templates\\insert_artist.html", nil)
	}
	if req.Method == "POST" {

		msg_artist := &artist.Artist{
			TableName:    "artist",
			Artist_ID:    1,
			Artist_Name:  req.PostFormValue("artist_name"),
			Artist_Email: req.PostFormValue("artist_email"),
		}
		if msg_artist.Validate() == false {
			render(w, "F:\\Client\\templates\\insert_artist.html", msg_artist)
			return
		}
		// send choose to server
		_, err = conn.Write([]byte("1"))
		if err != nil {
			log.Fatal(err.Error())
		}

		myslice := []string{msg_artist.TableName, msg_artist.Artist_Name, msg_artist.Artist_Email}
		data := strings.Join(myslice, ",")
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err.Error())
		}
		http.Redirect(w, req, "/", http.StatusSeeOther)

	}
}
func DeleteAlbum(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		//new
		_, err = conn.Write([]byte("2"))
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = conn.Write([]byte("album"))
		if err != nil {
			log.Fatal(err.Error())
		}
		buffer := make([]byte, 1024)
		data_size, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err.Error())
		}
		data_len := string(buffer[:data_size])
		len, _ := strconv.Atoi(data_len) // num of IDs of albums
		buffer1 := make([]byte, 1024)
		data_size, err = conn.Read(buffer1)
		if err != nil {
			log.Fatal(err.Error())
		}
		data := string(buffer1[:data_size])
		row := strings.Split(data, ";")
		var Ids []int
		for i := 0; i < len; i++ {
			id, _ := strconv.Atoi(row[i])
			Ids = append(Ids, id)
		}
		ids := IDS{
			IDs: Ids,
		}
		render(w, "F:\\Client\\templates\\delete_album.html", ids)
	}
	if req.Method == "POST" {
		id, _ := strconv.Atoi(req.PostFormValue("album_id"))
		str_ID := strconv.Itoa(id)
		_, err = conn.Write([]byte(str_ID))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
func DeleteArtist(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		//new
		_, err = conn.Write([]byte("2"))
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = conn.Write([]byte("artist"))
		if err != nil {
			log.Fatal(err.Error())
		}
		buffer := make([]byte, 1024)
		data_size, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err.Error())
		}
		data_len := string(buffer[:data_size])
		len, _ := strconv.Atoi(data_len) // num of IDs of artist
		buffer1 := make([]byte, 1024)
		data_size, err = conn.Read(buffer1)
		if err != nil {
			log.Fatal(err.Error())
		}
		data := string(buffer1[:data_size])
		row := strings.Split(data, ";")
		var Ids []int
		for i := 0; i < len; i++ {
			id, _ := strconv.Atoi(row[i])
			Ids = append(Ids, id)
		}
		ids := IDS{
			IDs: Ids,
		}
		render(w, "F:\\Client\\templates\\delete_artist.html", ids)
	}
	if req.Method == "POST" {
		id, _ := strconv.Atoi(req.PostFormValue("artist_id"))
		str_ID := strconv.Itoa(id)
		_, err = conn.Write([]byte(str_ID))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
func UpdateAlbum(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		render(w, "F:\\Client\\templates\\update_album.html", nil)
	}
	if req.Method == "POST" {
		id, _ := strconv.Atoi(req.PostFormValue("album_id"))
		artist_id, _ := strconv.Atoi(req.PostFormValue("artist_id"))
		price, _ := strconv.Atoi(req.PostFormValue("price"))

		msg_album := &tables.Album{
			TableName: "album",
			ID:        int64(id),
			Artist_ID: int64(artist_id),
			Title:     req.PostFormValue("title"),
			Price:     int64(price),
		}
		if msg_album.Validate() == false {
			render(w, "F:\\Client\\templates\\update_album.html", msg_album)
			return
		}
		_, err = conn.Write([]byte("3"))
		if err != nil {
			log.Fatal(err.Error())
		}
		str_ID := strconv.Itoa(id)
		str_Price := strconv.Itoa(price)
		str_artist_id := strconv.Itoa(artist_id)
		myslice := []string{msg_album.TableName, str_ID, msg_album.Title, str_artist_id, str_Price}
		data := strings.Join(myslice, ",")
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
func UpdateArtist(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		render(w, "F:\\Client\\templates\\update_artist.html", nil)
	}
	if req.Method == "POST" {
		artist_id, _ := strconv.Atoi(req.PostFormValue("artist_id"))
		msg_artist := &artist.Artist{
			TableName:    "artist",
			Artist_ID:    int64(artist_id),
			Artist_Name:  req.PostFormValue("artist_name"),
			Artist_Email: req.PostFormValue("artist_email"),
		}
		if msg_artist.Validate() == false {
			render(w, "F:\\Client\\templates\\update_artist.html", msg_artist)
			return
		}
		_, err = conn.Write([]byte("3"))
		if err != nil {
			log.Fatal(err.Error())
		}
		str_artist_id := strconv.Itoa(artist_id)
		myslice := []string{msg_artist.TableName, str_artist_id, msg_artist.Artist_Name, msg_artist.Artist_Email}
		data := strings.Join(myslice, ",")
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
func SelectAlbum(w http.ResponseWriter, req *http.Request) {
	_, err := conn.Write([]byte("4"))
	if err != nil {
		log.Fatal(err.Error())
	}
	Table_Name := "album"
	_, err = conn.Write([]byte(Table_Name))
	if err != nil {
		log.Fatal(err.Error())
	}
	buffer := make([]byte, 1024)
	data_size, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	data_len := string(buffer[:data_size])
	len, _ := strconv.Atoi(data_len) //num of Albums in table
	// fmt.Println(len)

	data_size, err = conn.Read(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := string(buffer[:data_size])
	// fmt.Println(data)
	row := strings.Split(data, ";")
	msg_album := &tables.Album{}
	albums := []tables.Album{}
	for i := 0; i < len; i++ {
		attributes := strings.Split(row[i], ",")
		id, _ := strconv.Atoi(attributes[0])
		artist_id, _ := strconv.Atoi(attributes[2])
		price, _ := strconv.Atoi(attributes[3])
		msg_album = &tables.Album{

			TableName: Table_Name,
			ID:        int64(id),
			Artist_ID: int64(artist_id),
			Title:     attributes[1],
			Price:     int64(price),
		}
		albums = append(albums, *msg_album)

	}
	Total_Album := Albums{
		AllAlbums: albums,
	}
	render(w, "F:\\Client\\templates\\select_album.html", Total_Album)
}
func SelectArtist(w http.ResponseWriter, req *http.Request) {
	_, err := conn.Write([]byte("4"))
	if err != nil {
		log.Fatal(err.Error())
	}
	Table_Name := "artist"
	_, err = conn.Write([]byte(Table_Name))
	if err != nil {
		log.Fatal(err.Error())
	}
	buffer := make([]byte, 1024)
	data_size, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	data_len := string(buffer[:data_size])
	len, _ := strconv.Atoi(data_len) //num of Artists in table
	// fmt.Println(len)

	data_size, err = conn.Read(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := string(buffer[:data_size])
	// fmt.Println(data)
	row := strings.Split(data, ";")
	artists := []artist.Artist{}
	msg_artist := &artist.Artist{}
	for i := 0; i < len; i++ {
		attributes := strings.Split(row[i], ",")

		artist_id, _ := strconv.Atoi(attributes[0])

		msg_artist = &artist.Artist{

			TableName:    Table_Name,
			Artist_ID:    int64(artist_id),
			Artist_Name:  attributes[1],
			Artist_Email: attributes[2],
		}
		artists = append(artists, *msg_artist)

	}
	Total_Artists := Artists{
		AllArtists: artists,
	}
	render(w, "F:\\Client\\templates\\select_artist.html", Total_Artists)
}
