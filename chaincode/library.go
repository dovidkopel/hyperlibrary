package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) Init() {

}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Book, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Book, error) {
	var assets []*Book
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Book
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
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
		inst := BookInstance{"bookInstance", instId, bookId, time.Now(), cost}
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

func (s *SmartContract) QueryBook(ctx contractapi.TransactionContextInterface, key string, value string) ([]*Book, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"book","%s":"%s"}}`, key, value)
	return getQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) BorrowBook(ctx contractapi.TransactionContextInterface, inst BookInstance, person Person) {

}

func (s *SmartContract) ReturnBook(ctx contractapi.TransactionContextInterface, inst BookInstance) {

}
