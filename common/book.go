package common

import (
	"time"
)

type Genre string
type Condition string
type Status string

const (
	FICTION     Genre = "FICTION"
	NON_FICTION       = "NON_FICTION"
	MYSTERY           = "MYSTERY"
)

const (
	NEW                  Condition = "NEW"
	GOOD                           = "GOOD"
	WORN                           = "WORN"
	RIPPED                         = "RIPPED"
	PAGES_MISSING                  = "PAGES_MISSING"
	REQUIRES_REPLACEMENT           = "REQUIRES_REPLACEMENT"
)

const (
	AVAILABLE Status = "AVAILABLE"
	RETURNED         = "RETURNED"
	RESERVED         = "RESERVED"
	OUT              = "OUT"
	LOST             = "LOST"
)

type Book struct {
	DocType   string `json:"docType" default:"book" `
	Isbn      string `json:"isbn"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Genre     Genre  `json:"genre"`
	Owned     uint   `json:"owned" default:"0"`
	Available uint   `json:"available" default:"0"`
	MaxId     uint16 `json:"maxId"`
}

type BookInstance struct {
	DocType   string    `json:"docType" default:"bookInstance"`
	Id        string    `json:"id"`
	BookId    string    `json:"bookId"`
	Purchased time.Time `json:"purchased"`
	Cost      float32   `json:"cost"`
	Status    Status    `json:"status"`
	Condition Condition `json:"condition"`
	DueDate   time.Time `json:"dueDate"`
	Borrower  User      `json:"borrower"`
}
