package helper

import "time"

type RequestStatus struct {
	Order_id string `json:"order_id" form:"order_id"`
}

type UpdateStatusTrx struct {
	Status string
}

type RequestUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type RequestForPayment struct {
	Bank     string `json:"bank_code"`
	Username string `json:"username"`
	Domain   string `json:"domain"`
}

type RequestCC struct {
	Card_number string `json:"card_number"`
	Cvv         string `json:"cvv"`
	Exp_month   int    `json:"exp_month"`
	Exp_year    int    `json:"exp_year"`
	Username    string `json:"username"`
	Domain      string `json:"domain"`
	// Otp st
}

type RequestUpdateCart struct {
	Status   int
	Order_id string
	Username string
	Domain   string
}

type RequestStoreTransaction struct {
	Order_id          string
	Username          string
	Date_insert       time.Time
	Date_payment      string
	Status            string
	Type_pembayaran   string
	TrxCode           string `gorm:"column:trxCode"`
	Unique_code       *int
	Jumlah            int
	Status_pembayaran *int
	Type_transaction  *string
	Mid_json_response string
	Mid_json_payment  *string
	Pg_type           int
	Xen_json_response *string
	Is_recurring      *int
	Installment_data  *string
	Is_mid_check      *int
	Card_number       *string
	Is_send_notif     *int
	Is_send_email     *int
	Email_to          string
}

type RequestCharge struct {
	Bank                string        `json:"bank"`
	Cart_id             []int         `json:"cart_id"`
	Payment_type        string        `json:"payment_type" form:"payment_type" query:"payment_type"`
	Transaction_details TrxDetails    `json:"transaction_details"`
	Customer_details    CsDetails     `json:"customer_details"`
	Item_details        []ItemDetails `json:"item_details"`
	Bank_transfer       BankTransfer  `json:"bank_transfer"`
}

type TrxDetails struct {
	Order_id     string `json:"order_id"`
	Gross_amount int    `json:"gross_amount"`
}

type CsDetails struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Phone      string `json:"phone"`
}

type ItemDetails struct {
	Id       string `json:"id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

type StoreTransaction struct {
	Order_id          string      `json:"order_id"`
	Username          string      `json:"username"`
	Date_insert       time.Time   `json:"date_insert"`
	Date_payment      time.Time   `json:"date_payment"`
	Status            string      `json:"status"`
	Type_pembayaran   string      `json:"type_pembayaran"`
	TrxCode           string      `json:"trx_code"`
	Unique_code       int         `json:"unique_code"`
	Jumlah            string      `json:"jumlah"`
	Status_pembayaran int         `json:"status_pembayaran"`
	Type_transaction  string      `json:"type_transaction"`
	Mid_json_response interface{} `json:"mid_json_response"`
	Mid_json_payment  interface{} `json:"mid_json_payment"`
	Pg_type           int         `json:"pg_type"`
	Xen_json_response interface{} `json:"xen_json_response"`
	Is_recurring      int         `json:"is_recurring"`
	Installment_data  string      `json:"installment_data"`
	Is_mid_check      int         `json:"is_mid_check"`
	Card_number       int         `json:"card_number"`
}
