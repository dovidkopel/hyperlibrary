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
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
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

func New(userId string, roles []string, handleEvents bool) LibraryClient {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists(userId) {
		//err = populateWallet(wallet)
		err = CreateAppUser(wallet, userId, roles)
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

	if handleEvents {
		ll.HandleEvents()
	}

	return ll
}

func blockEventHandler(c <-chan *fab.FilteredBlockEvent) {
	v := <-c

	txs := v.FilteredBlock.GetFilteredTransactions()
	for i := range txs {
		tx := txs[i]
		txa := tx.GetTransactionActions().GetChaincodeActions()
		for ii := range txa {
			act := txa[ii]
			ev := act.ChaincodeEvent

			var payload map[string]interface{}
			json.Unmarshal(ev.Payload, &payload)
			log.Println("EVENT", ev.EventName, payload)
		}
	}
}

func (l *LibraryClient) HandleEvents() {
	_, ch, _ := l.contract.RegisterEvent("Events")
	go l.eventHandler(ch)
	//_, ch, _ := l.network.RegisterFilteredBlockEvent()
	//go blockEventHandler(ch)
}

func (l *LibraryClient) eventHandler(c <-chan *fab.CCEvent) {
	v := <-c

	//event := v.EventName
	payloadBytes := v.Payload

	var payloads []common.Event
	err := json.Unmarshal(payloadBytes, &payloads)

	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, event := range payloads {
		pb, err := json.Marshal(event.Payload)

		if err != nil {
			log.Fatalf(err.Error())
		}

		if event.Name == "BookInstance.Created" {
			var inst common.BookInstance
			err = json.Unmarshal(pb, &inst)

			if err != nil {
				log.Fatalf(err.Error())
			}

			log.Println("EVENT", event.Name, inst)
		} else if event.Name == "BookInstance.Returned" {
			var inst common.BookInstance
			err = json.Unmarshal(pb, &inst)

			if err != nil {
				log.Fatalf(err.Error())
			}

			log.Println("EVENT", event.Name, inst)
			l.bookReturned(inst)
		} else {
			log.Println("EVENT", event.Name, event.Payload)
		}
	}
	go l.eventHandler(c)
}

func (l *LibraryClient) bookReturned(inst common.BookInstance) {

	r := rand.Intn(100)

	var cond common.Condition
	var fee float64 = 0
	available := true
	// Good

	if r > 50 {
		cond = common.GOOD
	} else if r > 40 {
		cond = common.WORN
	} else if r > 30 {
		cond = common.RIPPED
		fee = .50
	} else if r > 20 {
		cond = common.PAGES_MISSING
		fee = 1.0
	} else {
		cond = common.REQUIRES_REPLACEMENT
		fee = float64(inst.Cost)
		available = false
	}

	log.Println(fmt.Sprintf("Inspecting book with %s, %f", cond, fee))
	_, err := l.contract.SubmitTransaction("Invoke", "inspect", inst.Id, string(cond), fmt.Sprintf("%f", fee), strconv.FormatBool(available))

	if err != nil {
		log.Fatalf(err.Error())
	}

	if fee > 0 {
		log.Println(fee)
	}

	inst, _ = l.GetBookInstance(inst.Id)
	log.Println("inst", inst)
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

func (l *LibraryClient) ReturnBook(instId string) (common.Fee, error) {
	lateFeeBytes, err := l.contract.SubmitTransaction("ReturnBookInstance", instId)

	if err != nil {
		return common.Fee{}, err
	}

	var lateFee common.Fee
	err = json.Unmarshal(lateFeeBytes, &lateFee)

	if err != nil {
		return common.Fee{}, err
	}

	return lateFee, nil
}

func (l *LibraryClient) ListUsersOwingFees() ([]common.UserWithFees, error) {
	usersBytes, err := l.contract.EvaluateTransaction("ListUsersOwingFees")

	if err != nil {
		log.Fatalf(err.Error())
		return []common.UserWithFees{}, err
	}

	var users []common.UserWithFees
	err = json.Unmarshal(usersBytes, &users)

	if err != nil {
		return []common.UserWithFees{}, err
	}

	return users, nil
}

func (l *LibraryClient) PayLateFee(amount float64, feeIds []string) (common.Payment, error) {
	ids, err := json.Marshal(feeIds)

	if err != nil {
		return common.Payment{}, err
	}

	paymentBytes, err := l.contract.SubmitTransaction("Invoke", "pay", fmt.Sprintf("%f", amount), string(ids))

	if err != nil {
		return common.Payment{}, err
	}

	var payment common.Payment
	err = json.Unmarshal(paymentBytes, &payment)

	if err != nil {
		return common.Payment{}, err
	}

	return payment, nil
}

func (l *LibraryClient) GetFeeHistory(id string) ([]*common.History, error) {
	historyBytes, err := l.contract.EvaluateTransaction("GetFeeHistory", id)

	if err != nil {
		return nil, err
	}

	var history []*common.History
	err = json.Unmarshal(historyBytes, &history)

	if err != nil {
		return nil, err
	}

	return history, nil
}

func (l *LibraryClient) InspectReturnedBook(instId string, cond common.Condition, feeAmount float64, available bool) (*common.Fee, error) {
	feeBytes, err := l.contract.SubmitTransaction("InspectReturnedBook", instId, string(cond), fmt.Sprintf("%f", feeAmount), fmt.Sprintf("%d", available))

	if err != nil {
		return nil, err
	}

	var fee *common.Fee
	err = json.Unmarshal(feeBytes, &fee)

	if err != nil {
		return nil, err
	}

	return fee, nil
}
