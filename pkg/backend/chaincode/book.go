package chaincode

import (
	"time"
)

type Genre uint8
type Condition uint8
type Status uint8

const (
	FICTION Genre = iota
	NON_FICTION
)

const (
	NEW Condition = iota
	GOOD
	WORN
	RIPPED
	PAGES_MISSING
)

const (
	AVAILABLE Status = iota
	RESERVED
	OUT
	LOST
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
	DocType   string    `json:"docType" default:"bookInstance"`
	Id        string    `json:"id"`
	BookId    string    `json:"bookId"`
	Purchased time.Time `json:"purchased"`
	Cost      float32   `json:"cost"`
	Status    Status    `json:"status"`
	Condition Condition `json:"condition"`
}

type Person struct {
	Id        string
	FirstName string
	LastName  string
	phone     string
}
