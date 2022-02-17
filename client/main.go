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

	l := app.New("libraryApp@org1.example.com")
	//l1 := app.New("john@org1.example.com")
	//books := l.ListBooks()
	//for i := range books {
	//	b := books[i]
	//	log.Println(b)
	//}

	//l.CreateBook(common.Book{"book", "foobar5565554", "F. Scott Fitzgerald", "Dr Pepper", common.FICTION, 0, 0, 0})
	//l.PurchaseBook("foobar5565554", 5, 10.50)

	//bookInstances := l.ListBooksInstances("foobar5565554")
	//
	//var toTakeOut common.BookInstance = common.BookInstance{}
	//for i := range bookInstances {
	//	b := bookInstances[i]
	//	log.Println(b)
	//	if b.Status == common.AVAILABLE {
	//		toTakeOut = b
	//		break
	//	}
	//}
	//
	//err = l.BorrowBook(toTakeOut.Id)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//lateFee, err := l.ReturnBook(toTakeOut.Id)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//log.Println(fmt.Sprintf("Late Fee: %s", lateFee))

	//bookInstance, err := l.GetBookInstance("foobar5565554-1")
	//
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//log.Println("Book Instance", bookInstance)
	//
	//if bookInstance.Status == common.AVAILABLE {
	//	err = l.BorrowBook("foobar5565554-1`")
	//	if err != nil {
	//		log.Fatalf(err.Error())
	//	}
	//}
	//
	//lateFee, err := l.ReturnBook("foobar5565554-1")
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//log.Println(fmt.Sprintf("Late Fee: %s", lateFee))

	////
	//books := l.ListBooks()
	//for i := range books {
	//	b := books[i]
	//	log.Println(b)
	//}
	histories, _ := l.GetFeeHistory("aed2130379aa83e7102476315f3d7e58dbec65b05002f703fb957426ce6f2588")

	for i := range histories {
		history := histories[i]
		log.Println(history)
	}

	//users, err := l.ListUsersOwingFees()
	//for i := range users {
	//	user := users[i]
	//	log.Println("Users owing fees", user)
	//
	//	for k, _ := range user.FeesOwed {
	//		if k == "aed2130379aa83e7102476315f3d7e58dbec65b05002f703fb957426ce6f2588" {
	//			p, err := l.PayLateFee(1.0, []string{k})
	//
	//			if err != nil {
	//				log.Fatalf(err.Error())
	//			}
	//
	//			log.Println("Payment", p)
	//		}
	//	}
	//}
}
