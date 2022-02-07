package main

import (
	"hyperlibrary/client/app"
	"log"
	"os"
)

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	l := app.New("libraryApp6@org1.example.com")
	//print(l.ListBooks())
	books := l.ListBooks()
	for i := range books {
		b := books[i]
		log.Println(b)
	}

	//l.CreateBook(common.Book{"book", "foobar412443", "F. Scott Fitzgerald", "Blah1", common.FICTION, 0, 0, 0})
	l.PurchaseBook("abcd1234", 1, 10.50)

	books = l.ListBooks()
	for i := range books {
		b := books[i]
		log.Println(b)
	}
}
