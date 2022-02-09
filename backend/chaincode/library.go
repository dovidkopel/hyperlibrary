package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"hyperlibrary/common"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type SmartContract struct {
	contractapi.Contract
	BorrowDuration time.Duration
	LateFeePerDay  float64
}

func NewSmartContract(contract contractapi.Contract) *SmartContract {
	s := &SmartContract{Contract: contract}
	s.LateFeePerDay = .50
	bd, _ := time.ParseDuration(fmt.Sprintf("%dh", 14*24))
	s.BorrowDuration = bd
	return s
}

func (t *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	log.Println("Init invoked")

	t.CreateBook(ctx, common.Book{"book", "abcd1234", "Charles Dickens", "A Tale of Two Cities", common.FICTION, 0, 0, 0})
	t.CreateBook(ctx, common.Book{"book", "abcd45454", "William Shakespeare", "Romeo and Juliet", common.FICTION, 0, 0, 0})
	t.CreateBook(ctx, common.Book{"book", "abcd45455", "William Shakespeare", "Julis Casar", common.FICTION, 0, 0, 0})
	return nil
}

func (t *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) ([]byte, error) {
	log.Println("ex02 Invoke")
	if os.Getenv("DEVMODE_ENABLED") != "" {
		log.Println("invoking in devmode")
	}
	function, args := ctx.GetStub().GetFunctionAndParameters()
	name, _, _ := ctx.GetClientIdentity().GetAttributeValue("Name")
	clientId, _ := ctx.GetClientIdentity().GetID()
	log.Println(fmt.Sprintf("Client id: %s %s", clientId, name))

	switch args[0] {
	case "create":
		var book common.Book
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
	case "purchase":
		isbn := args[1]
		q, err := strconv.ParseUint(args[2], 10, 8)

		if err != nil {
			return nil, err
		}
		quantity := uint16(q)

		c, err := strconv.ParseFloat(args[3], 32)

		if err != nil {
			return nil, err
		}
		cost := float32(c)

		insts, err := t.PurchaseBook(ctx, isbn, quantity, cost)
		print("foo")

		if err != nil {
			return nil, err
		}

		_, err = json.Marshal(insts)

		if err != nil {
			return nil, err
		}

		//log.Println(fmt.Sprintf("Purchase ret: %s", string(ret)))

		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf(`Invalid invoke "%s" function name. Expecting "invoke", "delete", "query", "respond", "mspid", or "event"`, function))
	}
}

func (t *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, book common.Book) error {
	assetBytes, err := json.Marshal(book)
	if err != nil {
		return err
	}

	// Check for existing ISBN
	books, err := t.QueryBook(ctx, "isbn", book.Isbn)
	if len(books) > 0 {
		return errors.New(fmt.Sprintf(`A book with the "%s" ISBN already exists!`, book.Isbn))
	}

	ctx.GetStub().SetEvent("Book.Created", assetBytes)
	return ctx.GetStub().PutState("book."+book.Isbn, assetBytes)
}

func (t *SmartContract) PurchaseBook(ctx contractapi.TransactionContextInterface, bookId string, quantity uint16, cost float32) ([]*common.BookInstance, error) {
	if quantity < 1 {
		return nil, errors.New("Quantity must be at least 1!")
	}

	assetBytes, err := ctx.GetStub().GetState("book." + bookId)

	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", bookId, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", bookId)
	}

	var book common.Book
	err = json.Unmarshal(assetBytes, &book)

	log.Println(fmt.Sprintf("There are currently %d owned.", book.Owned))

	var instances []*common.BookInstance

	var i uint16
	starting_id := book.MaxId + 1
	var last_id uint16

	for i = 0; i <= quantity; i++ {
		instId := fmt.Sprintf("%s-%d", book.Isbn, starting_id+i)
		last_id = starting_id + i

		inst := common.BookInstance{"bookInstance", instId, bookId, cost,
			common.AVAILABLE, common.NEW, time.Time{}, common.User{}}
		instBytes, err := json.Marshal(inst)

		if err != nil {
			log.Println("Unable to marshal instance!")
			return nil, err
		}

		err = ctx.GetStub().PutState("bookInstance."+instId, instBytes)
		ctx.GetStub().SetEvent("BookInstance.Created", instBytes)

		if err != nil {
			log.Println("Unable to store instance state!")
			return nil, err
		}

		instances = append(instances, &inst)
	}

	book.Owned += uint(quantity)
	book.Available += uint(quantity)
	book.MaxId = last_id

	assetBytes, err = json.Marshal(book)
	if err != nil {
		log.Println("Unable to marshal book!")
		return nil, err
	}

	err = ctx.GetStub().PutState("book."+bookId, assetBytes)

	log.Println(fmt.Sprintf("Created %d instances", quantity))
	//return instances, nil
	return nil, nil
}

func (t *SmartContract) QueryBook(ctx contractapi.TransactionContextInterface, key string, value string) ([]*common.Book, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"book","%s":"%s"}}`, key, value)
	res, err := getQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, err
	}

	var books []*common.Book
	for i := range res {
		bookBytes := res[i]
		var book common.Book
		err = json.Unmarshal(bookBytes, &book)
		books = append(books, &book)
	}
	return books, nil
}

func (t *SmartContract) ListBooks(ctx contractapi.TransactionContextInterface) ([]*common.Book, error) {
	res, err := getQueryResultForQueryString(ctx, `{"selector":{"docType":"book"}}`)

	if err != nil {
		return nil, err
	}

	var books []*common.Book
	for i := range res {
		bookBytes := res[i]
		var book common.Book
		err = json.Unmarshal(bookBytes, &book)
		books = append(books, &book)
	}
	return books, nil
}

func (t *SmartContract) ListBookInstances(ctx contractapi.TransactionContextInterface, isbn string, statuses []common.Status) ([]*common.BookInstance, error) {
	selector := map[string]interface{}{
		"docType": "bookInstance",
		"bookId":  isbn,
	}

	if len(statuses) > 0 {
		var orStatuses []map[string]common.Status
		for i := range statuses {
			orStatuses = append(orStatuses, map[string]common.Status{
				"status": statuses[i],
			})
		}
		selector["$or"] = orStatuses
	} else {

	}

	query := map[string]interface{}{
		"selector": selector,
	}

	queryString, err := json.Marshal(query)

	if err != nil {
		return nil, err
	}

	res, err := getQueryResultForQueryString(ctx, string(queryString))

	if err != nil {
		return nil, err
	}

	var books []*common.BookInstance
	for i := range res {
		bookBytes := res[i]
		var book common.BookInstance
		err = json.Unmarshal(bookBytes, &book)
		books = append(books, &book)
	}
	return books, nil
}

func (t *SmartContract) GetBook(ctx contractapi.TransactionContextInterface, isbn string) (*common.Book, error) {
	bookBytes, err := ctx.GetStub().GetState(fmt.Sprintf("book.%s", isbn))

	if err != nil {
		return nil, err
	}

	var book common.Book
	err = json.Unmarshal(bookBytes, &book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (t *SmartContract) GetBookInstance(ctx contractapi.TransactionContextInterface, instId string) (*common.BookInstance, error) {
	bookBytes, err := ctx.GetStub().GetState(fmt.Sprintf("bookInstance.%s", instId))

	if err != nil {
		return nil, err
	}

	var bookInstance common.BookInstance
	err = json.Unmarshal(bookBytes, &bookInstance)

	if err != nil {
		return nil, err
	}

	return &bookInstance, nil
}

func (t *SmartContract) UpdateBook(ctx contractapi.TransactionContextInterface, book *common.Book) error {
	bookBytes, err := json.Marshal(book)

	if err != nil {
		return err
	}

	ctx.GetStub().PutState(fmt.Sprintf("book.%s", book.Isbn), bookBytes)
	return nil
}

func (t *SmartContract) BorrowBookInstance(ctx contractapi.TransactionContextInterface, instId string) error {
	instBytes, err := ctx.GetStub().GetState(fmt.Sprintf("bookInstance.%s", instId))

	if err != nil {
		return err
	}

	var inst common.BookInstance
	err = json.Unmarshal(instBytes, &inst)

	if err != nil {
		return err
	}

	book, err := t.GetBook(ctx, inst.BookId)

	if inst.Status == common.AVAILABLE {
		log.Println(fmt.Sprintf("Going to borrow book \"%s\"", instId))
		clientId, _ := ctx.GetClientIdentity().GetID()
		name, _, _ := ctx.GetClientIdentity().GetAttributeValue("Name")

		inst.Borrower = common.User{clientId, name}
		inst.Status = common.OUT

		//inst.DueDate = time.Now().Add(t.BorrowDuration )
		dd, _ := time.ParseDuration(fmt.Sprintf("-%dh", 5*24))
		inst.DueDate = time.Now().Add(dd).Round(time.Hour)
		instBytes1, err := json.Marshal(&inst)

		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(fmt.Sprintf("bookInstance.%s", instId), instBytes1)

		if err != nil {
			return err
		}

		book.Available -= 1
		err = t.UpdateBook(ctx, book)

		if err != nil {
			return err
		}

		ctx.GetStub().SetEvent("BookInstance.Borrowed", instBytes1)

		return nil
	} else if inst.Status == common.OUT {
		return errors.New("This book is already out!")
	}

	return nil
}

func (t *SmartContract) ReturnBookInstance(ctx contractapi.TransactionContextInterface, instId string) (common.LateFee, error) {
	instBytes, err := ctx.GetStub().GetState(fmt.Sprintf("bookInstance.%s", instId))

	if err != nil {
		return common.LateFee{}, err
	}

	var inst common.BookInstance
	err = json.Unmarshal(instBytes, &inst)

	if err != nil {
		return common.LateFee{}, err
	}

	book, err := t.GetBook(ctx, inst.BookId)

	if err != nil {
		return common.LateFee{}, err
	}

	if inst.Status == common.OUT {
		now := time.Now()

		// Late fee
		if inst.DueDate.Before(now) {
			diff := now.Sub(inst.DueDate).Round(time.Hour)
			diffDays := math.RoundToEven(diff.Hours() / 24)

			if diffDays > 0 {
				log.Println("A late fee is owed")
				id := ctx.GetStub().GetTxID()
				fee := t.LateFeePerDay * diffDays
				ts, _ := ctx.GetStub().GetTxTimestamp()
				date := time.Unix(ts.Seconds, int64(ts.Nanos)).Round(time.Hour).UTC()

				lateFee := common.LateFee{id, inst.Borrower, fee, date, false}
				log.Println(lateFee)
				lateFeeBytes, err := json.Marshal(lateFee)

				if err != nil {
					return common.LateFee{}, err
				}

				err = ctx.GetStub().PutState(fmt.Sprintf("lateFee.%s", id), lateFeeBytes)

				if err != nil {
					return common.LateFee{}, err
				}

				err = ctx.GetStub().SetEvent("LateFee.Created", lateFeeBytes)

				if err != nil {
					return common.LateFee{}, err
				}

				return lateFee, nil
			}
		}

		inst.Status = common.AVAILABLE
		inst.DueDate = time.Time{}
		inst.Borrower = common.User{}

		instBytes, err = json.Marshal(inst)

		if err != nil {
			return common.LateFee{}, err
		}

		ctx.GetStub().PutState(fmt.Sprintf("bookInstance.%s", instId), instBytes)
		ctx.GetStub().SetEvent("BookInstance.Returned", instBytes)

		book.Available += 1
		t.UpdateBook(ctx, book)
	} else {
		errors.New("Book cannot be returned if it is not out!")
	}

	return common.LateFee{}, err
}
