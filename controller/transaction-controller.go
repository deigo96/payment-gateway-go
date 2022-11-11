package controller

import (
	"encoding/json"
	"net/http"
	"topup-service/config"
	"topup-service/helper"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func (controller *AuthController) CheckTransaction(c echo.Context) error {
	req := new(helper.RequestStatus)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.Order_id == "" {
		resp := helper.BuildErrorResponse("Order id is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	var status int
	r := coreapi.Client{}
	r.New(config.GetConfig().Server_key, midtrans.EnvironmentType(config.GetConfig().Midtrans_env))

	transactionStatusResp, e := r.CheckTransaction(req.Order_id)
	if e != nil {
		resp := helper.BuildErrorResponse(e.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				out, _ := json.Marshal(transactionStatusResp)
				json_payment := string(out)
				status = 2

				update := helper.RequestStoreTransaction{
					Status:           transactionStatusResp.TransactionStatus,
					Mid_json_payment: &json_payment,
				}
				err := controller.transaction.UpdateStoreTransactionService(update, transactionStatusResp.OrderID, status)
				if err != nil {
					return err
				}
				// TODO set transaction status on your databaase to 'success'
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				out, _ := json.Marshal(transactionStatusResp)
				json_payment := string(out)
				status = 3

				update := helper.RequestStoreTransaction{
					Status:           transactionStatusResp.TransactionStatus,
					Mid_json_payment: &json_payment,
				}
				err := controller.transaction.UpdateStoreTransactionService(update, transactionStatusResp.OrderID, status)
				if err != nil {
					return err
				}
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				status = 1
				out, _ := json.Marshal(transactionStatusResp)
				json_payment := string(out)

				update := helper.RequestStoreTransaction{
					Status:           transactionStatusResp.TransactionStatus,
					Mid_json_payment: &json_payment,
				}
				err := controller.transaction.UpdateStoreTransactionService(update, transactionStatusResp.OrderID, status)
				if err != nil {
					return err
				}
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}

	resp := helper.BuildSuccessResponse(true, "Succes", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}
