package client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"log"
)

type LibraryClient struct {
}

func (l *LibraryClient) Init() *gateway.Contract {
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromRaw()),
		gateway.WithIdentity()
	)

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	return network.GetContract("hyperlibrary")
}

func (l *LibraryClient) CreateBook() {
	l.Init()

}
