package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"hyperlibrary/common"
	"log"
	"math"
	"time"
)

func (t *SmartContract) GetLateFees(ctx contractapi.TransactionContextInterface, clientId string) ([]*common.LateFee, error) {
	res, err := GetQueryResultForQueryString(ctx, fmt.Sprintf(`{"selector":{"docType":"lateFee", "borrower.clientId": "%s"}}`, clientId))

	if err != nil {
		return nil, err
	}

	var lateFees []*common.LateFee
	for i := range res {
		lateFeeBytes := res[i]
		var lateFee common.LateFee
		err = json.Unmarshal(lateFeeBytes, &lateFee)
		lateFees = append(lateFees, &lateFee)
	}
	return lateFees, nil
}

func (t *SmartContract) CheckForLateFee(ctx contractapi.TransactionContextInterface, inst common.BookInstance) (common.LateFee, error) {
	now := time.Now()

	if inst.DueDate.Before(now) {
		diff := now.Sub(inst.DueDate).Round(time.Hour)
		diffDays := math.RoundToEven(diff.Hours() / 24)

		if diffDays > 0 {
			log.Println("A late fee is owed")
			id := ctx.GetStub().GetTxID()
			fee := t.LateFeePerDay * diffDays
			ts, _ := ctx.GetStub().GetTxTimestamp()
			date := common.GetApproxTime(ts)

			lateFee := common.LateFee{"lateFee", id, inst.Borrower, fee, date, 0.0, false}

			lateFeeBytes, err := t.StoreFee(ctx, lateFee)

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

	return common.LateFee{}, nil
}

func (t *SmartContract) GetFee(ctx contractapi.TransactionContextInterface, feeId string) (common.LateFee, error) {
	feeBytes, err := ctx.GetStub().GetState(fmt.Sprintf("lateFee.%s", feeId))

	if err != nil {
		return common.LateFee{}, err
	}

	var fee common.LateFee
	err = json.Unmarshal(feeBytes, &fee)

	if err != nil {
		return common.LateFee{}, err
	}

	return fee, nil
}

func (t *SmartContract) StoreFee(ctx contractapi.TransactionContextInterface, fee common.LateFee) ([]byte, error) {
	lateFeeBytes, err := json.Marshal(fee)

	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(fmt.Sprintf("lateFee.%s", fee.Id), lateFeeBytes)

	if err != nil {
		return nil, err
	}

	return lateFeeBytes, nil
}

func (t *SmartContract) DistributePayment(ctx contractapi.TransactionContextInterface, amount float64, feeIds []string) (map[string]float64, error) {
	remainingAmount := amount
	fees := make(map[string]float64, len(feeIds))

	for i := range feeIds {
		feeId := feeIds[i]
		fee, err := t.GetFee(ctx, feeId)

		if err != nil {
			return map[string]float64{}, err
		}

		if fee.FullyPaid {
			continue
		}

		remainingFee := fee.Fee - fee.AmountPaid

		if remainingAmount > 0 {
			if remainingAmount >= remainingFee {
				fee.AmountPaid += remainingFee
				remainingAmount -= remainingFee
				fees[feeId] = remainingFee
			} else {
				ableToPay := remainingFee - remainingAmount
				fee.AmountPaid += ableToPay
				remainingAmount -= ableToPay
				fees[feeId] = ableToPay
			}

			var feeEvent string
			if fee.AmountPaid == fee.Fee {
				fee.FullyPaid = true
				feeEvent = "LateFee.FullyPaid"
			} else {
				feeEvent = "LateFee.PartiallyPaid"
			}

			feeBytes, err := t.StoreFee(ctx, fee)

			if err != nil {
				return map[string]float64{}, err
			}

			ctx.GetStub().SetEvent(feeEvent, feeBytes)
		} else {
			fees[feeId] = 0
			log.Println("For the fee there isn't enough money left in the payment.", fee)
		}
	}

	return fees, nil
}

func (t *SmartContract) StorePayment(ctx contractapi.TransactionContextInterface, payment common.Payment) (common.Payment, error) {
	paymentBytes, err := json.Marshal(payment)

	if err != nil {
		log.Fatalf(err.Error())
		return common.Payment{}, err
	}

	err = ctx.GetStub().PutState(fmt.Sprintf("payment.%s", payment.Id), paymentBytes)

	if err != nil {
		log.Fatalf(err.Error())
		return common.Payment{}, err
	}

	err = ctx.GetStub().SetEvent("Payment.Created", paymentBytes)

	if err != nil {
		log.Fatalf(err.Error())
		return common.Payment{}, err
	}

	return payment, nil
}

func (t *SmartContract) PayLateFee(ctx contractapi.TransactionContextInterface, amount float64, feeIds []string) (common.Payment, error) {
	feesPaid, err := t.DistributePayment(ctx, amount, feeIds)

	log.Println("Fees to be paid", feesPaid)

	if err != nil {
		return common.Payment{}, err
	}

	if len(feesPaid) > 0 {
		txId := ctx.GetStub().GetTxID()
		ts, _ := ctx.GetStub().GetTxTimestamp()
		date := common.GetApproxTime(ts)

		payment := common.Payment{"payment", txId,
			t.GetUserByClientId(ctx), amount, date, feesPaid,
		}

		log.Println("Creating payment", payment)

		return t.StorePayment(ctx, payment)
	}

	return common.Payment{}, nil
}

func (t *SmartContract) GetFeeHistory(ctx contractapi.TransactionContextInterface, id string) ([]common.History, error) {
	return t.GetHistory(ctx, fmt.Sprintf("lateFee.%s", id))
}
