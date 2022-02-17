package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"hyperlibrary/common"
	"log"
)

func (t *SmartContract) GetUserByClientId(ctx contractapi.TransactionContextInterface) common.User {
	clientId, _ := ctx.GetClientIdentity().GetID()
	name, _, _ := ctx.GetClientIdentity().GetAttributeValue("Name")
	return common.User{clientId, name}
}

func (t *SmartContract) ListUsersOwingFees(ctx contractapi.TransactionContextInterface) ([]common.UserWithFees, error) {
	res, err := GetQueryResultForQueryString(ctx, `{"selector":{"docType":"lateFee","fullyPaid":false}}`)
	var users []common.UserWithFees

	if err != nil {
		return users, err
	}

	log.Println("Late fees", len(res))

	for i := range res {
		lateFeeBytes := res[i]
		var lateFee common.LateFee
		err = json.Unmarshal(lateFeeBytes, &lateFee)

		if err != nil {
			return users, err
		}

		log.Println("Late Fee", lateFee)

		found := false
		for i := range users {
			user := users[i]
			if user.ClientId == lateFee.Borrower.ClientId {
				found = true
				am := lateFee.Fee - lateFee.AmountPaid
				user.FeesOwed[lateFee.Id] = am
				user.TotalOwed += am
				users[i] = user
				break
			}
		}

		if !found {
			am := lateFee.Fee - lateFee.AmountPaid
			fo := map[string]float64{lateFee.Id: am}
			user := common.UserWithFees{
				User:      lateFee.Borrower,
				FeesOwed:  fo,
				TotalOwed: am,
			}
			users = append(users, user)
		}
	}

	log.Println("Users owing fees", users)

	return users, nil
}
