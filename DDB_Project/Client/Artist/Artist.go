package artist

import (
	"regexp"
	"strings"
)

var rxEmail = regexp.MustCompile(".+@.+\\.+")

type Artist struct {
	TableName    string
	Artist_ID    int64
	Artist_Name  string
	Artist_Email string
	Errors       map[string]string
}

func (msg *Artist) Validate() bool {
	msg.Errors = make(map[string]string)
	match := rxEmail.Match([]byte(msg.Artist_Email))
	if match == false {
		msg.Errors["Artist_Email"] = "Please enter valied email"
	}
	if msg.TableName != "artist" {
		msg.Errors["TableName"] = "Error Table Name"
	}
	if strings.TrimSpace(msg.Artist_Name) == "" {
		msg.Errors["Artist_Name"] = "Please enter the title"
	}
	if msg.Artist_ID <= 0 {
		msg.Errors["Artist_ID"] = "Enter positive number"
	}

	return len(msg.Errors) == 0
}
