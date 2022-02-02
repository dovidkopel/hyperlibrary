package main

import (
	"github.com/dovidkopel/hyperlibrary/backend/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	libraryChaincode, err := contractapi.NewChaincode(chaincode.NewSmartContract(contractapi.Contract{}))
	if err != nil {
		log.Panicf("Error creating library chaincode: %v", err)
	}

	if err := libraryChaincode.Start(); err != nil {
		log.Panicf("Error starting library chaincode: %v", err)
	}
}
