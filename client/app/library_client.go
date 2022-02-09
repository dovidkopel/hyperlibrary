package app

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"hyperlibrary/common"
	"log"
	"os"
	"path/filepath"
)

type LibraryClient struct {
	network  *gateway.Network
	gateway  *gateway.Gateway
	contract *gateway.Contract
}

var fabric_samples = "/home/dkopel/go/src/github.com/dovidkopel/fabric-samples"

func getConfig() core.ConfigProvider {
	os.Setenv("FABRIC_SDK_GO_PROJECT_PATH", fabric_samples)
	//os.Setenv("CRYPTOCONFIG_FIXTURES_PATH", "test-network/organizations/cryptogen")
	//ccpPath := filepath.Join(
	//	"config.yaml",
	//)
	ccpPath := "/home/dkopel/go/src/hyperlibrary/client/config.yaml"
	return config.FromFile(filepath.Clean(ccpPath))
}

func New(userId string) LibraryClient {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists(userId) {
		//err = populateWallet(wallet)
		err = CreateAppUser(wallet, userId)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(getConnectionConfig()),
		gateway.WithIdentity(wallet, userId),
	)

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	ll := LibraryClient{}
	ll.network = network
	ll.gateway = gw
	ll.contract = network.GetContract("hyperlibrary")
	ll.HandleEvents()

	return ll
}

func eventHandler(c <-chan *fab.FilteredBlockEvent) {
	v := <-c
	log.Println(v)
}

func (l *LibraryClient) HandleEvents() {
	_, _, _ = l.contract.RegisterEvent("Book.Created")
	_, _, _ = l.contract.RegisterEvent("BookInstance.Created")
	_, _, _ = l.contract.RegisterEvent("BookInstance.Borrowed")
	_, _, _ = l.contract.RegisterEvent("BookInstance.Returned")
	_, ch, _ := l.network.RegisterFilteredBlockEvent()
	go eventHandler(ch)
}

func (l *LibraryClient) ListBooks() []common.Book {
	print("Listing books")
	resp, err := l.contract.EvaluateTransaction("ListBooks")

	var books []common.Book
	json.Unmarshal(resp, &books)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return books
}

func (l *LibraryClient) ListBooksInstances(isbn string) []common.BookInstance {
	print("Listing book instances")
	resp, err := l.contract.EvaluateTransaction("ListBookInstances", isbn, "[]")

	if err != nil {
		log.Fatalf(err.Error())
	}

	var books []common.BookInstance
	err = json.Unmarshal(resp, &books)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return books
}

func (l *LibraryClient) CreateBook(book common.Book) error {
	payload, err := json.Marshal(book)
	_, err = l.contract.SubmitTransaction("Invoke", "create", string(payload))

	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	return nil
}

func (l *LibraryClient) PurchaseBook(isbn string, quantity int, cost float32) ([]common.BookInstance, error) {
	instBytes, err := l.contract.SubmitTransaction("Invoke", "purchase", isbn,
		fmt.Sprintf("%d", quantity),
		fmt.Sprintf("%f", cost),
	)

	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}

	var insts []common.BookInstance
	json.Unmarshal(instBytes, &insts)

	return insts, nil
}

func (l *LibraryClient) GetBookInstance(instId string) (common.BookInstance, error) {
	bookInstanceBytes, err := l.contract.SubmitTransaction("GetBookInstance", instId)

	if err != nil {
		return common.BookInstance{}, err
	}

	var bookInstance common.BookInstance
	err = json.Unmarshal(bookInstanceBytes, &bookInstance)

	if err != nil {
		return common.BookInstance{}, err
	}

	return bookInstance, err
}

func (l *LibraryClient) BorrowBook(bookId string) error {
	_, err := l.contract.SubmitTransaction("BorrowBookInstance", bookId)

	if err != nil {
		return err
	}
	return nil
}

func (l *LibraryClient) ReturnBook(instId string) (common.LateFee, error) {
	lateFeeBytes, err := l.contract.SubmitTransaction("ReturnBookInstance", instId)

	if err != nil {
		return common.LateFee{}, err
	}

	var lateFee common.LateFee
	err = json.Unmarshal(lateFeeBytes, &lateFee)

	if err != nil {
		return common.LateFee{}, err
	}

	return lateFee, nil
}
