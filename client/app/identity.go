package app

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"log"
	"path/filepath"
)

func getConnectionConfig() core.ConfigProvider {
	ccpPath := filepath.Join(
		fabric_samples,
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)
	return config.FromFile(filepath.Clean(ccpPath))
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		fabric_samples,
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		//"Admin@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

func CreateAppUser(wallet *gateway.Wallet, id string) error {
	sdk, err := fabsdk.New(getConfig())
	print("created sdk")
	if err != nil {
		log.Fatalf(err.Error())
	}

	mspClient, err := mspclient.New(
		sdk.Context(),
		mspclient.WithOrg("Org1"),
		mspclient.WithCAInstance("ca.org1.example.com"),
	)

	if err != nil {
		log.Fatalf(err.Error())
	}

	var attrs []mspclient.Attribute

	attrs = append(attrs,
		mspclient.Attribute{"Name", id, true},
		mspclient.Attribute{"Role", "library", true},
	)

	secret, err := mspClient.Register(&mspclient.RegistrationRequest{
		Name:        id,
		Type:        "client",
		Affiliation: "org1.department1",
		Attributes:  attrs,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	err = mspClient.Enroll(id,
		mspclient.WithSecret(secret),
		mspclient.WithProfile("tls"),
		mspclient.WithType("app"),
	)

	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	resp, err := mspClient.GetSigningIdentity(id)

	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	if resp != nil {
		key, _ := resp.PrivateKey().Bytes()
		StoreWallet(wallet, id, string(resp.EnrollmentCertificate()), string(key))
	}

	return nil
}

func StoreWallet(wallet *gateway.Wallet, id string, cert string, key string) {
	identity := gateway.NewX509Identity("Org1MSP", cert, key)
	wallet.Put(id, identity)
}
