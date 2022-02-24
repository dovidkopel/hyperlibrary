package main

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	assert2 "github.com/stretchr/testify/assert"
	"hyperlibrary/client/app"
	"hyperlibrary/common"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano() + time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func TestMain(m *testing.M) {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	files, _ := os.ReadDir("wallet")
	for _, f := range files {
		os.Remove(fmt.Sprintf("wallet/%s", f.Name()))
	}

	code := m.Run()
	os.Exit(code)
}

func Test1(t *testing.T) {
	memberName := fmt.Sprintf(`member-%s@org1.example.com`, RandomString(10))
	librarianName := fmt.Sprintf(`librarian-%s@org1.example.com`, RandomString(10))
	member := app.New(memberName, []string{"MEMBER"}, false)
	librarian := app.New(librarianName, []string{"LIBRARIAN"}, false)

	librarian.HandleEvents()

	book1 := common.Book{"book", RandomString(18), RandomString(10), RandomString(10), common.FICTION, 0, 0, 0}
	book2 := common.Book{"book", RandomString(18), RandomString(10), RandomString(10), common.FICTION, 0, 0, 0}
	book3 := common.Book{"book", RandomString(18), RandomString(10), RandomString(10), common.FICTION, 0, 0, 0}

	err := librarian.CreateBook(&book1)
	assert.Equal(t, err, nil)

	err = librarian.CreateBook(&book2)
	assert.Equal(t, err, nil)

	err = librarian.CreateBook(&book3)
	assert.Equal(t, err, nil)

	insts, err := librarian.PurchaseBook(book1.Isbn, 10, 5.50)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(insts), 10)

	// Check instances
	for _, inst := range insts {
		assert.Equal(t, inst.Status, common.AVAILABLE)
		assert.Equal(t, inst.Condition, common.NEW)
		assert.Equal(t, inst.DueDate, time.Time{})
		assert.Equal(t, inst.Cost, float32(5.5))
		assert.Equal(t, inst.Borrower, common.User{Roles: []string{}})
	}

	// Make sure the counts match
	insts1 := member.ListBooksInstances(book1.Isbn, []common.Status{common.AVAILABLE})
	assert.Equal(t, len(insts1), 10)

	out, _ := member.GetMyBooksOut()
	assert.Equal(t, len(out), 0)

	// Borrow one book of book1
	borrowedInst, err := member.BorrowBookInstance(insts1[0].Id)
	assert.Equal(t, err, nil)

	out, _ = member.GetMyBooksOut()
	assert.Equal(t, len(out), 1)
	assert.Equal(t, string(borrowedInst.Status), common.OUT)
	assert.Equal(t, borrowedInst.Borrower.Name, memberName)

	// Try to take the same book again
	borrowedInst1, err := member.BorrowBookInstance(insts[0].Id)
	assert2.True(t, borrowedInst1 == nil)
	assert2.True(t, err != nil)

	// Make sure that there are now only 9 books available of book1
	insts2 := member.ListBooksInstances(book1.Isbn, []common.Status{common.AVAILABLE})
	assert.Equal(t, len(insts2), 9)

	instsA, err := librarian.PurchaseBook(book2.Isbn, 10, 50.25)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(instsA), 10)

	// Borrow one book of book2
	borrowedInstA, err := member.BorrowBookInstance(instsA[0].Id)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(borrowedInstA.Status), common.OUT)
	assert.Equal(t, borrowedInstA.Borrower.Name, memberName)

	instsA2 := librarian.ListBooksInstances(book2.Isbn, []common.Status{common.AVAILABLE})
	assert.Equal(t, len(instsA2), 9)

	instsB, err := librarian.PurchaseBook(book3.Isbn, 10, 20.00)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(instsB), 10)

	// Try to take out more books than the max
	_, err = member.BorrowBookInstance(instsB[0].Id)
	assert2.True(t, err != nil)

	myFees, err := member.GetMyUnpaidFees()
	assert.Equal(t, len(myFees), 0)

	librarian.RegisterEventHandler("BookInstance.Returned", func(pb []byte) {
		var inst *common.BookInstance
		err := json.Unmarshal(pb, &inst)

		if err != nil {
			log.Fatalf(err.Error())
		}

		log.Println("EVENT", inst)
		fee, err := librarian.InspectReturnedBook(inst.Id, common.WORN, .25, true)
		assert.Equal(t, fee.Fee, float64(.25))

		// After
		memberFees, err := member.GetMyUnpaidFees()
		assert.Equal(t, len(memberFees), 2)

		amt := float64(0)
		var ids []string
		for _, fee := range memberFees {
			amt += fee.Fee
			ids = append(ids, fee.Id)
		}
		payment, err := member.PayFee(amt, ids)
		assert.Equal(t, payment.Amount, amt)

		memberFees, err = member.GetMyUnpaidFees()
		assert.Equal(t, len(memberFees), 0)
	})

	fee, err := member.ReturnBookInstance(borrowedInst.Id)
	assert.Equal(t, fee.Fee, float64(5.5))
}
