package main

import (
	"hyperlibrary/client/app"
	"hyperlibrary/common"
	"log"
	"os"
)

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	l := app.New("libraryApp@org1.example.com")
	books := l.ListBooks()
	for i := range books {
		b := books[i]
		log.Println(b)
	}

	//l.CreateBook(common.Book{"book", "foobar412443", "F. Scott Fitzgerald", "Blah1", common.FICTION, 0, 0, 0})
	//l.PurchaseBook("foobar412443", 2, 10.50)

	bookInstances := l.ListBooksInstances("foobar412443")

	var toTakeOut common.BookInstance = common.BookInstance{}
	for i := range bookInstances {
		b := bookInstances[i]

		if b.Status == common.AVAILABLE {
			toTakeOut = b
			break
		}
	}

	l.BorrowBook(toTakeOut.Id)

	books = l.ListBooks()
	for i := range books {
		b := books[i]
		log.Println(b)
	}
}
