package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"topup-service/config"
	handler "topup-service/handler"
	"topup-service/helper"
	"topup-service/helper/api"
	"topup-service/helper/cart"
	"topup-service/helper/payment"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type AuthController struct {
	jwtService  handler.JwtService
	cart        cart.CartService
	transaction payment.TransactionService
}

func NewAuthController(jwtService handler.JwtService, cart cart.CartService, transaction payment.TransactionService) *AuthController {
	return &AuthController{
		jwtService:  jwtService,
		cart:        cart,
		transaction: transaction,
	}
}

func (controller *AuthController) GetListBank(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	if header == "Bearer " || header == "" {
		resp := helper.BuildErrorResponse("No Authorization, No Token", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	validate := controller.jwtService.ValidateToken(header)
	if validate != true {
		resp := helper.BuildErrorResponse("Token expired or invalid", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	VA := []helper.ListBank{
		{Bank_code: "BCA", Bank_name: "BCA"},
		{Bank_code: "PER", Bank_name: "PERMATA"},
		{Bank_code: "BNI", Bank_name: "BNI"},
		{Bank_code: "ATM", Bank_name: "Jaringan ATM"},
	}

	Instant := []helper.ListBank{
		{Bank_code: "GOPAY", Bank_name: "GoPay"},
		{Bank_code: "QRIS", Bank_name: "QRIS"},
		{Bank_code: "CC", Bank_name: "Kartu Kredit"},
		{Bank_code: "MBP", Bank_name: "Mandiri Bill Payment"},
	}

	Convenience := []helper.ListBank{
		{Bank_code: "ALFA", Bank_name: "Alfamart / Alfamidi / Dan+Dan"},
		{Bank_code: "INDO", Bank_name: "Indomaret"},
	}

	category := map[string][]helper.ListBank{
		"Transfer Bank / Virtual Account": VA,
		"Instant Payment":                 Instant,
		"Convenience Store":               Convenience,
	}

	resp := helper.BuildSuccessResponse(true, "Success get bank list", category)
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) PaymentGateway(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	if header == "Bearer " || header == "" {
		resp := helper.BuildErrorResponse("No Authorization, No Token", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	validate := controller.jwtService.ValidateToken(header)
	if validate != true {
		resp := helper.BuildErrorResponse("Token expired or invalid", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	r := coreapi.Client{}
	r.New(config.GetConfig().Server_key, midtrans.EnvironmentType(config.GetConfig().Midtrans_env))

	t := time.Now()
	timeString := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
	order := t.UnixMilli()
	orderId := "INV-001" + strconv.FormatInt(order, 10)

	requestr := new(helper.RequestForPayment)
	if err := c.Bind(requestr); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	getUser := api.GetUserProfile(requestr.Username)
	bank := strings.ToUpper(requestr.Bank)
	chargeReq := &coreapi.ChargeReq{}
	item := midtrans.ItemDetails{}
	ga := 0
	sliceitem := []midtrans.ItemDetails{}
	cartData, _ := controller.cart.GetCartService(requestr.Username, requestr.Domain)
	if len(cartData.Results) < 1 {
		resp := helper.BuildErrorResponse("Your cart is empty", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}
	for _, val := range cartData.Results {
		ga += val.Harga
		item.ID = strconv.Itoa(val.Top_up_pack_id)
		item.Name = val.Information
		item.Qty = 1
		item.Price = int64(val.Harga)
		sliceitem = append(sliceitem, item)
	}

	switch bank {
	case "BCA":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("bank_transfer"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.Bank("bca"),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
			CustomExpiry: &coreapi.CustomExpiry{
				OrderTime:      timeString,
				ExpiryDuration: 2,
				Unit:           "minute",
			},
		}
	case "BNI":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("bank_transfer"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.Bank("bni"),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
		}
	case "PER", "ATM":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("permata"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
		}
	case "MBP":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("echannel"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			EChannel: &coreapi.EChannelDetail{
				BillInfo1: "Payment:",
				BillInfo2: "Online purchase",
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
		}
	case "ALFA":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("cstore"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
			ConvStore: &coreapi.ConvStoreDetails{
				Store:   "alfamart",
				Message: "GPS.id Top Up Package",
			},
		}
	case "INDO":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("cstore"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
			ConvStore: &coreapi.ConvStoreDetails{
				Store:   "indomaret",
				Message: "GPS.id Top Up Package",
			},
		}
	case "QRIS":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("qris"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
			Qris: &coreapi.QrisDetails{
				Acquirer: "gopay",
			},
		}
	case "GOPAY":
		chargeReq = &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType("gopay"),
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: int64(ga),
			},
			Items: &sliceitem,
			CustomerDetails: &midtrans.CustomerDetails{
				Email: getUser.Email,
				LName: getUser.Username,
				Phone: getUser.Phone,
			},
		}
	default:
		chargeReq = &coreapi.ChargeReq{}
	}

	coreApiRes, _ := r.ChargeTransaction(chargeReq)
	if coreApiRes.StatusCode == "400" {
		a := map[string]interface{}{
			"status_code":    coreApiRes.StatusCode,
			"status_message": coreApiRes.StatusMessage,
		}

		resp := helper.BuildErrorResponse(coreApiRes.StatusMessage, a)
		return c.JSON(http.StatusOK, resp)
	}

	checkPreviousTrx, _ := controller.transaction.CheckPreviousTrxService(requestr.Username)
	if len(checkPreviousTrx) > 0 {
		for _, val := range checkPreviousTrx {
			_, _ = r.CancelTransaction(val.Order_id)
			_ = controller.cart.UpdateStatusCartPaymentService(val.Order_id)

		}
	}

	update := helper.RequestUpdateCart{
		Status:   1,
		Order_id: orderId,
		Username: requestr.Username,
		Domain:   requestr.Domain,
	}
	_ = controller.cart.UpdateCartPaymentService(update)
	gross, _ := strconv.ParseFloat(coreApiRes.GrossAmount, 64)
	out, err := json.Marshal(coreApiRes)
	if err != nil {
		panic(err)
	}

	store := helper.RequestStoreTransaction{
		Order_id:        orderId,
		Username:        getUser.Username,
		Date_insert:     time.Now(),
		Date_payment:    coreApiRes.TransactionTime,
		Status:          coreApiRes.TransactionStatus,
		Type_pembayaran: coreApiRes.PaymentType,
		TrxCode:         coreApiRes.TransactionID,
		// Unique_code      : ,
		Jumlah: int(gross),
		// Status_pembayaran: ,
		// Type_transaction : ,
		Mid_json_response: string(out),
		// Pg_type          : ,
		// Xen_json_response: ,
		// Is_recurring     : ,
		// Installment_data : ,
		// Is_mid_check     : ,
		// Card_number      : ,
		// Is_send_notif    : ,
		// Is_send_email    : ,
		Email_to: getUser.Email,
	}
	// fmt.Println(store)

	_ = controller.transaction.StoreTransactionService(store)
	dataRes := payment.ResponseVa(coreApiRes, bank)
	m := dataRes.(map[string]interface{})

	response := helper.BuildSuccessResponse(true, m["status_message"].(string), dataRes)
	return c.JSON(http.StatusOK, response)

}

func (controller *AuthController) PaymentCreditCartd(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	if header == "Bearer " || header == "" {
		resp := helper.BuildErrorResponse("No Authorization, No Token", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	validate := controller.jwtService.ValidateToken(header)
	if validate != true {
		resp := helper.BuildErrorResponse("Token expired or invalid", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	r := coreapi.Client{}
	r.New(config.GetConfig().Server_key, midtrans.EnvironmentType(config.GetConfig().Midtrans_env))

	t := time.Now()
	timeString := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
	order := t.UnixMilli()
	orderId := "INV-001" + strconv.FormatInt(order, 10)

	req := new(helper.RequestCC)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if req.Username == "" {
		resp := helper.BuildErrorResponse("Username is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Domain == "" {
		resp := helper.BuildErrorResponse("Domain is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Card_number == "" {
		resp := helper.BuildErrorResponse("Card number is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Exp_month < 1 || req.Exp_month > 12 {
		resp := helper.BuildErrorResponse("Invalid expired month", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Exp_year < time.Now().Year() {
		resp := helper.BuildErrorResponse("Invalid expired year", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	getUser := api.GetUserProfile(req.Username)
	if getUser.Id == 0 {
		resp := helper.BuildErrorResponse("User not found", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	cartData, _ := controller.cart.GetCartService(req.Username, req.Domain)
	if len(cartData.Results) < 1 {
		resp := helper.BuildErrorResponse("Your cart is empty", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	clientKey := config.GetConfig().Client_key

	getToken, _ := r.CardToken(req.Card_number, req.Exp_month, req.Exp_year, req.Cvv, clientKey)
	if getToken.StatusCode != "200" {
		resp := helper.BuildErrorResponse(getToken.ValidationMessage[0], helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	item := midtrans.ItemDetails{}
	ga := 0
	sliceitem := []midtrans.ItemDetails{}
	for _, val := range cartData.Results {
		ga += val.Harga
		item.ID = strconv.Itoa(val.Top_up_pack_id)
		item.Name = val.Information
		item.Qty = 1
		item.Price = int64(val.Harga)
		sliceitem = append(sliceitem, item)
	}

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType("credit_card"),
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: int64(ga),
		},
		CreditCard: &coreapi.CreditCardDetails{
			TokenID:        getToken.TokenID,
			Authentication: true,
		},
		Items: &sliceitem,
		CustomerDetails: &midtrans.CustomerDetails{
			Email: getUser.Email,
			LName: getUser.Username,
			Phone: getUser.Phone,
		},
		CustomExpiry: &coreapi.CustomExpiry{
			OrderTime:      timeString,
			ExpiryDuration: 2,
			Unit:           "minute",
		},
	}

	coreApiRes, _ := r.ChargeTransaction(chargeReq)
	if coreApiRes.StatusCode != "201" {
		a := map[string]interface{}{
			"status_code":    coreApiRes.StatusCode,
			"status_message": coreApiRes.StatusMessage,
		}

		resp := helper.BuildErrorResponse(coreApiRes.StatusMessage, a)
		return c.JSON(http.StatusOK, resp)
	}

	checkPreviousTrx, _ := controller.transaction.CheckPreviousTrxService(req.Username)
	if len(checkPreviousTrx) > 0 {
		for _, val := range checkPreviousTrx {
			_, _ = r.CancelTransaction(val.Order_id)
			_ = controller.cart.UpdateStatusCartPaymentService(val.Order_id)

		}
	}

	update := helper.RequestUpdateCart{
		Status:   1,
		Order_id: orderId,
		Username: req.Username,
		Domain:   req.Domain,
	}
	_ = controller.cart.UpdateCartPaymentService(update)
	gross, _ := strconv.ParseFloat(coreApiRes.GrossAmount, 64)
	out, err := json.Marshal(coreApiRes)
	if err != nil {
		panic(err)
	}

	store := helper.RequestStoreTransaction{
		Order_id:        orderId,
		Username:        getUser.Username,
		Date_insert:     time.Now(),
		Date_payment:    coreApiRes.TransactionTime,
		Status:          coreApiRes.TransactionStatus,
		Type_pembayaran: coreApiRes.PaymentType,
		TrxCode:         coreApiRes.TransactionID,
		// Unique_code      : ,
		Jumlah: int(gross),
		// Status_pembayaran: ,
		// Type_transaction : ,
		Mid_json_response: string(out),
		// Pg_type          : ,
		// Xen_json_response: ,
		// Is_recurring     : ,
		// Installment_data : ,
		// Is_mid_check     : ,
		// Card_number      : ,
		// Is_send_notif    : ,
		// Is_send_email    : ,
		Email_to: getUser.Email,
	}

	_ = controller.transaction.StoreTransactionService(store)
	dataRes := payment.ResponseVa(coreApiRes, "CC")
	m := dataRes.(map[string]interface{})

	resp := helper.BuildSuccessResponse(true, m["status_message"].(string), dataRes)
	return c.JSON(http.StatusOK, resp)
}

// func (controller *AuthController) CheckTransaction(c echo.Context) (err error) {
// 	request := new(helper.RequestStatus)
// 	if err = c.Bind(request); err != nil {
// 		return c.String(http.StatusBadRequest, "bad request")
// 	}

// 	r := coreapi.Client{}
// 	r.New(config.GetConfig().Server_key, midtrans.EnvironmentType(config.GetConfig().Midtrans_env))

// 	// configEnv := helper.ConfigAppEnv{
// 	// 	Server_key:   config.GetConfig().Server_key,
// 	// 	Midtrans_url: config.GetConfig().Midtrans_url,
// 	// }
// 	// a, e := r.
// 	// get, err := payment.GetTransactionStatus(request.Order_id, configEnv)
// 	get, e := r.CheckTransaction(request.Order_id)
// 	if e != nil {
// 		resp := helper.BuildErrorResponse(e.GetMessage(), helper.EmptyObj{})
// 		return c.JSON(http.StatusOK, resp)
// 	}

// 	response := helper.BuildSuccessResponse(true, "Success get transaction status", get)
// 	return c.JSON(http.StatusOK, response)
// }
