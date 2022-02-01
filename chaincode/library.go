package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"os"
	"time"
)

// CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["init"]}' --isInit

// CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=true CORE_CHAINCODE_ID_NAME=hyperlibrary:1.0 ./hyperlibrary -peer.address 127.0.0.1:7052

//  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fabcar --peerAddresses localhost:7051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --isInit -c '{"function":"initLedger","Args":[]}'

type SmartContract struct {
	contractapi.Contract
}

func (t *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Init invoked")

	t.CreateBook(ctx, Book{"book", "abcd1234", "Charles Dickens", "A Tale of Two Cities", FICTION, 0, 0})
	t.CreateBook(ctx, Book{"book", "abcd45454", "William Shakespeare", "Romeo and Juliet", FICTION, 0, 0})
	return nil
}

func (t *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) ([]byte, error) {
	fmt.Println("ex02 Invoke")
	if os.Getenv("DEVMODE_ENABLED") != "" {
		fmt.Println("invoking in devmode")
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
