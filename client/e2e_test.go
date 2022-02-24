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

	code := m.Run()
	os.Exit(code)
}

func Test1(t *testing.T) {
	memberName := fmt.Sprintf(`member-%s@org1.example.com`, RandomString(8))
	librarian := app.New(fmt.Sprintf(`librarian-%s@org1.example.com`, RandomString(8)), []string{"LIBRARIAN"}, false)

	librarian.HandleEvents()

	member := app.New(memberName, []string{"MEMBER"}, false)
	book := common.Book{"book", RandomString(18), RandomString(10), RandomString(10), common.FICTION, 0, 0, 0}
	err := librarian.CreateBook(&book)
	assert.Equal(t, err, nil)
	insts, err := librarian.PurchaseBook(book.Isbn, 10, 5.50)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(insts), 10)

	insts1 := librarian.ListBooksInstances(book.Isbn, []common.Status{common.AVAILABLE})
	assert.Equal(t, len(insts1), 10)

	for _, inst := range insts {
		assert.Equal(t, inst.Status, common.AVAILABLE)
		assert.Equal(t, inst.Condition, common.NEW)
		assert.Equal(t, inst.DueDate, time.Time{})
		assert.Equal(t, inst.Cost, float32(5.5))
		assert.Equal(t, inst.Borrower, common.User{Roles: []string{}})
	}

	borrowedInst, err := member.BorrowBookInstance(insts[0].Id)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(borrowedInst.Status), common.OUT)
	assert.Equal(t, borrowedInst.Borrower.Name, memberName)

	insts2 := librarian.ListBooksInstances(book.Isbn, []common.Status{common.AVAILABLE})
	assert.Equal(t, len(insts2), 9)

	// Try to take the same book again
	borrowedInst1, err := member.BorrowBookInstance(insts[0].Id)
	assert2.True(t, borrowedInst1 == nil)
	assert2.True(t, err != nil)

	myFees, err := member.GetMyFees()
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
		memberFees, err := member.GetMyFees()
		assert.Equal(t, len(memberFees), 2)
	})

	fee, err := member.ReturnBookInstance(borrowedInst.Id)
	assert.Equal(t, fee.Fee, float64(5.5))

}
