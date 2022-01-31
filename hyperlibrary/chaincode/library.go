package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, book Book) error {
	assetBytes, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("book."+book.Isbn, assetBytes)
}

func (s *SmartContract) PurchaseBook(ctx contractapi.TransactionContextInterface, bookId string, quantity uint8, cost float32) ([]BookInstance, error) {
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
		inst := BookInstance{instId, bookId, time.Now(), cost}
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

func (s *SmartContract) QueryBook(ctx contractapi.TransactionContextInterface, key string, value string) []Book {
	return []Book{}
}

func (s *SmartContract) BorrowBook(ctx contractapi.TransactionContextInterface, inst BookInstance, person Person) {

}

func (s *SmartContract) ReturnBook(ctx contractapi.TransactionContextInterface, inst BookInstance) {

}
