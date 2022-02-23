package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"hyperlibrary/backend/chaincode"

	"log"
)

func main() {
	contract := chaincode.NewSmartContract(contractapi.Contract{})
	contract.TransactionContextHandler = new(contractapi.TransactionContext)
	//contract.AfterTransaction = contract.After
	//contract.BeforeTransaction = contract.Before

	libraryChaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		log.Panicf("Error creating library chaincode: %v", err)
	}

	if err := libraryChaincode.Start(); err != nil {
		log.Panicf("Error starting library chaincode: %v", err)
	}
}
