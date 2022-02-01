package chaincode

import (
	"time"
)

type Genre uint8

const (
	FICTION     Genre = 0
	NON_FICTION Genre = 1
)

type Book struct {
	DocType   string `json:"docType" default:"book"`
	Isbn      string `json:"isbn"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Genre     Genre  `json:"genre"`
	Owned     uint   `json:"owned" default:"0"`
	Available uint   `json:"available" default:"0"`
}

type BookInstance struct {
	DocType   string `json:"docType" default:"bookInstance"`
	Id        string
	BookId    string
	Purchased time.Time
	Cost      float32
}

type Person struct {
	Id        string
	FirstName string
	LastName  string
	phone     string
}
