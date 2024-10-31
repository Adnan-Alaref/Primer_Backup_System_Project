package tables

import (
	"strings"
)

type Album struct {
	TableName string
	ID        int64
	Artist_ID int64
	Title     string
	Price     int64
	Errors    map[string]string
}

// type Data struct {
// 	albums []Album
// }

func (msg *Album) Validate() bool {
	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Title) == "" {
		msg.Errors["Title"] = "Please enter the title"
	}

	if msg.Price <= 0 {
		msg.Errors["Price"] = "Enter positive number"
	}

	return len(msg.Errors) == 0
}
