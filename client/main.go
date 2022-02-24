package main

import (
	"fmt"
	"hyperlibrary/client/app"
	"hyperlibrary/common"
	"log"
	"os"
	"time"
)

func borrowAndReturn(l app.LibraryClient) {
	bookInstances := l.ListBooksInstances("foobar5565554", []common.Status{})

	var toTakeOut common.BookInstance = common.BookInstance{}
	for i := range bookInstances {
		b := bookInstances[i]
		log.Println(b)
		if b.Status == common.AVAILABLE {
			toTakeOut = b
			break
		}
	}

	_, err := l.BorrowBookInstance(toTakeOut.Id)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("Borrowed..")

	lateFee, err := l.ReturnBookInstance(toTakeOut.Id)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println(fmt.Sprintf("Late Fee: %f", lateFee.Fee))

	//time.Sleep(10 * time.Second)

	log.Println(l.PayLateFee(lateFee.Fee, []string{lateFee.Id}))
}

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	user := "libraryApp@org1.example.com"
	roles := []string{"MEMBER"}

	cmd := ""
	for i := range os.Args {
		if os.Args[i] == "1" {
			cmd = "1"
			user = "john@org1.example.com"
			roles = append(roles, "LIBRARIAN")
		} else if os.Args[i] == "1.5" {
			cmd = "1.5"
			user = "john@org1.example.com"
			roles = append(roles, "LIBRARIAN")
		} else if os.Args[i] == "2" {
			cmd = "2"
			user = "tony@org1.example.com"
		} else if os.Args[i] == "3" {
			cmd = "3"
			user = "tom@org1.example.com"
			roles = append(roles, "LIBRARIAN")
		}
	}

	switch cmd {
	case "1":
		l := app.New(user, roles, false)
		l.CreateBook(&common.Book{"book", "foobar5565554", "F. Scott Fitzgerald", "Dr Pepper", common.FICTION, 0, 0, 0})
		l.PurchaseBook("foobar5565554", 5, 10.50)

		l.CreateBook(&common.Book{"book", "abdbd55687", "F. Scott Fitzgerald", "Root Beer", common.FICTION, 0, 0, 0})
		l.PurchaseBook("abdbd55687", 10, 50.0)
		break
	case "1.5":
		l := app.New(user, roles, false)
		insts, _ := l.PurchaseBook("foobar5565554", 5, 10.50)

		if len(insts) != 5 {
			panic("Expecting 5")
		}

		insts, _ = l.PurchaseBook("abdbd55687", 10, 50.0)

		if len(insts) != 10 {
			panic("Expecting 10")
		}
		break
	case "2":
		l := app.New(user, roles, false)
		borrowAndReturn(l)
		break
	case "3":
		l := app.New(user, roles, false)
		borrowAndReturn(l)
		break
	default:
		roles = append(roles, "LIBRARIAN")
		user = "sally@org1.example.com"
		app.New(user, roles, true)

		for true {
			//log.Println("Waiting 5 seconds...")
			time.Sleep(5 * time.Second)

			//users, err := l.ListUsersOwingFees()
			//
			//if err != nil {
			//	log.Fatalf(err.Error())
			//}
			//
			//for i := range users {
			//	user := users[i]
			//	log.Println("Users owing fees", user)
			//}
		}
	}
	return

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
	//histories, _ := l.GetFeeHistory("aed2130379aa83e7102476315f3d7e58dbec65b05002f703fb957426ce6f2588")
	//
	//for i := range histories {
	//	history := histories[i]
	//	log.Println(history)
	//}

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
