package payment

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"topup-service/helper"

	"github.com/midtrans/midtrans-go/coreapi"
)

type TransactionList interface {
	StoreTransactionList(helper.RequestStoreTransaction) error
	CheckPreviousTrxList(username string) ([]helper.RequestStatus, error)
	UpdateStoreTransactionList(helper.RequestStoreTransaction, string, int) error
}

type TransactionService interface {
	StoreTransactionService(helper.RequestStoreTransaction) error
	CheckPreviousTrxService(username string) ([]helper.RequestStatus, error)
	UpdateStoreTransactionService(helper.RequestStoreTransaction, string, int) error
}

type transactionService struct {
	transaction TransactionList
}

func NewTransactionService(transaction TransactionList) TransactionService {
	return &transactionService{
		transaction: transaction,
	}
}

func (t *transactionService) StoreTransactionService(s helper.RequestStoreTransaction) error {
	err := t.transaction.StoreTransactionList(s)
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionService) CheckPreviousTrxService(username string) (s []helper.RequestStatus, err error) {
	res, err := t.transaction.CheckPreviousTrxList(username)
	if err != nil {
		return res, nil
	}
	return res, nil
}

func (t *transactionService) UpdateStoreTransactionService(s helper.RequestStoreTransaction, orderId string, status int) error {
	err := t.transaction.UpdateStoreTransactionList(s, orderId, status)
	if err != nil {
		return err
	}
	return nil
}

func GetTransactionStatus(orderId string, d helper.ConfigAppEnv) (res interface{}, err error) {

	sandboxUrl := d.Midtrans_url
	serverKey := d.Server_key
	serverKey += ":"
	key := base64.StdEncoding.EncodeToString([]byte(serverKey))
	url := fmt.Sprintf(sandboxUrl + "/v2/" + orderId + "/status")

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("Error while reading the response bytes:", err)
	// }

	json.NewDecoder(resp.Body).Decode(&res)

	// str := string(body)
	// parser, _ := gojq.NewStringQuery(str)
	// fmt.Println(body)
	// status_code, _ := parser.QueryToString("status_code")
	// status_message, _ := parser.QueryToString("status_message")

	// if status_code == "404" {
	// 	return res, errors.New(status_message)
	// }

	// currency, _ := parser.QueryToString("currency")
	// fraud_status, _ := parser.QueryToString("fraud_status")
	// gross_amount, _ := parser.QueryToString("gross_amount")
	// merchant_id, _ := parser.QueryToString("merchant_id")
	// order_id, _ := parser.QueryToString("order_id")
	// payment_type, _ := parser.QueryToString("payment_type")
	// signature_key, _ := parser.QueryToString("signature_key")
	// transaction_id, _ := parser.QueryToString("transaction_id")
	// transaction_status, _ := parser.QueryToString("transaction_status")
	// transaction_time, _ := parser.QueryToString("transaction_time")
	// payment_amounts, _ := parser.QueryToArray("payment_amounts")
	// permata_va_number, _ := parser.QueryToString("permata_va_number")
	// va_numbers, _ := parser.QueryToArray("va_numbers")

	// var s helper.Status
	// s.Currency = currency
	// s.Fraud_status = fraud_status
	// s.Gross_amount = gross_amount
	// s.Merchant_id = merchant_id
	// s.Order_id = order_id
	// s.Status_code = status_code
	// s.Status_message = status_message
	// s.Payment_type = payment_type
	// s.Signature_key = signature_key
	// s.Transaction_id = transaction_id
	// s.Transaction_status = transaction_status
	// s.Transaction_time = transaction_time
	// s.Payment_amounts = payment_amounts
	// s.Va_numbers = va_numbers
	// s.Permata_va_number = permata_va_number

	return res, err
}

func ApproveTransaction(orderId string) (res interface{}, err error) {

	sandboxUrl := os.Getenv("SANDBOX_URL")
	serverKey := os.Getenv("SERVER_KEY")
	serverKey += ":"
	key := base64.StdEncoding.EncodeToString([]byte(serverKey))
	url := fmt.Sprintf(sandboxUrl + "/v2/" + orderId + "/approve")
	// fmt.Println(serverKey)

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("Error while reading the response bytes:", err)
	// }

	// str := string(body)
	// parser, _ := gojq.NewStringQuery(str)
	// fmt.Println(str)
	var resps map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&resps)

	return resp, err
}

func ResponseVa(c *coreapi.ChargeResponse, bank string) interface{} {
	switch bank {
	case "BCA":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"va_numbers":         c.VaNumbers,
			"fraud_status":       c.FraudStatus,
		}

		return res

	case "BNI":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"va_numbers":         c.VaNumbers,
			"fraud_status":       c.FraudStatus,
		}
		return res

	case "PER", "ATM":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"permata_va_number":  c.PermataVaNumber,
			"fraud_status":       c.FraudStatus,
		}
		return res

	case "MBP":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"fraud_status":       c.FraudStatus,
			"bill_key":           c.BillKey,
			"biller_code":        c.BillerCode,
		}
		return res

	case "ALFA":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"fraud_status":       c.FraudStatus,
			"payment_code":       c.PaymentCode,
			"store":              c.Store,
		}
		return res

	case "INDO":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"payment_code":       c.PaymentCode,
			"store":              c.Store,
			// "merchant_id":        c.FraudStatus,
		}
		return res

	case "QRIS":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"actions":            c.Actions,
			"qr_string":          c.QRString,
			"acquirer":           c.Acquirer,
			// "merchant_id":        c.FraudStatus,
		}
		return res

	case "GOPAY":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"actions":            c.Actions,
			// "merchant_id":        c.FraudStatus,
		}
		return res

	case "CC":

		res := map[string]interface{}{
			"status_code":        c.StatusCode,
			"status_message":     c.StatusMessage,
			"bank":               c.Bank,
			"transaction_id":     c.TransactionID,
			"order_id":           c.OrderID,
			"redirect_url":       c.RedirectURL,
			"fraud_status":       c.FraudStatus,
			"masked_card":        c.MaskedCard,
			"card_type":          c.CardType,
			"gross_amount":       c.GrossAmount,
			"currency":           c.Currency,
			"payment_type":       c.PaymentType,
			"transaction_time":   c.TransactionTime,
			"transaction_status": c.TransactionStatus,
			"on_us":              c.OnUs,
			// "merchant_id":        c.FraudStatus,
		}
		return res
	}

	return nil

}
