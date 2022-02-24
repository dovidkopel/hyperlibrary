package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"hyperlibrary/common"
	"log"
	"strings"
)

func (t *SmartContract) GetUserByClientId(ctx contractapi.TransactionContextInterface) common.User {
	clientId, _ := ctx.GetClientIdentity().GetID()
	name, _, _ := ctx.GetClientIdentity().GetAttributeValue("Name")
	roles, _, _ := ctx.GetClientIdentity().GetAttributeValue("Roles")

	rs := strings.Split(roles, ",")
	return common.User{clientId, name, rs}
}

func (t *SmartContract) HasRole(ctx contractapi.TransactionContextInterface, role string) bool {
	user := t.GetUserByClientId(ctx)

	for _, r := range user.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (t *SmartContract) ListUsersOwingFees(ctx contractapi.TransactionContextInterface) ([]common.UserWithFees, error) {
	res, err := GetQueryResultForQueryString(ctx, `{"selector":{"docType":"fee","fullyPaid":false}}`)
	var users []common.UserWithFees

	if err != nil {
		return users, err
	}

	log.Println("Fees", len(res))

	for i := range res {
		feeBytes := res[i]
		var fee common.Fee
		err = json.Unmarshal(feeBytes, &fee)

		if err != nil {
			return users, err
		}

		log.Println("Fee", fee)

		found := false
		for i := range users {
			user := users[i]
			if user.ClientId == fee.Borrower.ClientId {
				found = true
				am := fee.Fee - fee.AmountPaid
				user.FeesOwed[fee.Id] = am
				user.TotalOwed += am
				users[i] = user
				break
			}
		}

		if !found {
			am := fee.Fee - fee.AmountPaid
			fo := map[string]float64{fee.Id: am}
			user := common.UserWithFees{
				User:      fee.Borrower,
				FeesOwed:  fo,
				TotalOwed: am,
			}
			users = append(users, user)
		}
	}

	log.Println("Users owing fees", users)

	return users, nil
}
