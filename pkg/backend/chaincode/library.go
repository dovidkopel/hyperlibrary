package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"os"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

func NewSmartContract(contract contractapi.Contract) *SmartContract {
	s := &SmartContract{Contract: contract}
	return s
}

func (t *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	log.Println("Init invoked")

	t.CreateBook(ctx, Book{"book", "abcd1234", "Charles Dickens", "A Tale of Two Cities", FICTION, 0, 0})
	t.CreateBook(ctx, Book{"book", "abcd45454", "William Shakespeare", "Romeo and Juliet", FICTION, 0, 0})
	t.CreateBook(ctx, Book{"book", "abcd45455", "William Shakespeare", "Julis Casar", FICTION, 0, 0})
	return nil
}

func (t *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) ([]byte, error) {
	log.Println("ex02 Invoke")
	if os.Getenv("DEVMODE_ENABLED") != "" {
		log.Println("invoking in devmode")
	}
	function, args := ctx.GetStub().GetFunctionAndParameters()
	switch args[0] {
	case "create":
		var book Book
		err := json.Unmarshal([]byte(args[1]), &book)

		if err != nil {
			return nil, err
		}

		book.DocType = "book"
		book.Owned = 0
		book.Available = 0

		err = t.CreateBook(ctx, book)

		if err != nil {
			return nil, err
		}

		return nil, nil
	case "list":
		// Deletes an entity from its state
		books, err := t.ListBooks(ctx)

		if err != nil {
			return nil, err
		}

		ret, err := json.Marshal(books)

		if err != nil {
			return nil, err
		}

		return ret, nil
	default:
		return nil, errors.New(fmt.Sprintf(`Invalid invoke "%s" function name. Expecting "invoke", "delete", "query", "respond", "mspid", or "event"`, function))
	}
}

func (t *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, book Book) error {
	assetBytes, err := json.Marshal(book)
	if err != nil {
		return err
	}

	// Check for existing ISBN
	books, err := t.QueryBook(ctx, "isbn", book.Isbn)
	if len(books) > 0 {
		return errors.New(fmt.Sprintf(`A book with the "%s" ISBN already exists!`, book.Isbn))
	}

	return ctx.GetStub().PutState("book."+book.Isbn, assetBytes)
}

func (t *SmartContract) PurchaseBook(ctx contractapi.TransactionContextInterface, bookId string, quantity uint8, cost float32) ([]BookInstance, error) {
	assetBytes, err := ctx.GetStub().GetState(bookId)

	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", bookId, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", bookId)
	}

	var book Book
	err = json.Unmarshal(assetBytes, &book)
	instances := []BookInstance{}

	var i uint8
	for i = 0; i <= quantity; i++ {
		instId := uuid.New().String()
		inst := BookInstance{"bookInstance", instId, bookId, time.Now(), cost,
			AVAILABLE, NEW}
		instBytes, err := json.Marshal(inst)

		if err != nil {
			return instances, err
		}

		err = ctx.GetStub().PutState("bookInstance."+instId, instBytes)
		if err != nil {
			return instances, err
		}

		instances = append(instances, inst)
	}

	book.Owned += uint(quantity)
	book.Available += uint(quantity)

	assetBytes, err = json.Marshal(book)
	if err != nil {
		return instances, err
	}

	err = ctx.GetStub().PutState("book."+bookId, assetBytes)

	return instances, nil
}

func (t *SmartContract) QueryBook(ctx contractapi.TransactionContextInterface, key string, value string) ([]*Book, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"book","%s":"%s"}}`, key, value)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *SmartContract) ListBooks(ctx contractapi.TransactionContextInterface) ([]*Book, error) {
	return getQueryResultForQueryString(ctx, `{"selector":{"docType":"book"}}`)
}

func (t *SmartContract) BorrowBook(ctx contractapi.TransactionContextInterface, inst BookInstance, person Person) {

}

func (t *SmartContract) ReturnBook(ctx contractapi.TransactionContextInterface, inst BookInstance) {

}
