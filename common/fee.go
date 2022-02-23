package common

import "time"

type FeeType string

const (
	LATE_FEE   FeeType = "LATE"
	DAMAGE_FEE         = "DAMAGE"
	LOST_FEE           = "LOST"
)

type Fee struct {
	DocType    string    `json:"docType" default:"fee"`
	Id         string    `json:"id"`
	Borrower   User      `json:"borrower"`
	Fee        float64   `json:"fee"`
	Type       FeeType   `json:"feeType"`
	Date       time.Time `json:"date"`
	AmountPaid float64   `json:"amountPaid"`
	FullyPaid  bool      `json:"fullyPaid"`
}

type Payment struct {
	DocType string             `json:"docType" default:"payment"`
	Id      string             `json:"id"`
	Payer   User               `json:"payer"`
	Amount  float64            `json:"amount"`
	Date    time.Time          `json:"date"`
	FeeIds  map[string]float64 `json:"feeIds"`
}
