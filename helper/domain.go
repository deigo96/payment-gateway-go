package helper

type ExtractToken struct {
	Status  bool
	Message string
	Data    DataToken
}

type DataToken struct {
	Id       int
	Email    string
	Username string
	Phone    string
}

type AllData struct {
	TotalAllData int
	Category     []Category
}

type Category struct {
	Bank_category string `json:"bank_category"`
	Bank_list     []ListBank
}

type ListBank struct {
	Bank_code string `json:"bank_code"`
	Bank_name string `json:"bank_name"`
}

type Status struct {
	Currency           string      `json:"currency"`
	Fraud_status       string      `json:"fraud_status"`
	Gross_amount       string      `json:"gross_amount"`
	Merchant_id        string      `json:"merchant_id"`
	Order_id           string      `json:"order_id"`
	Payment_amounts    interface{} `json:"payment_amounts"`
	Payment_type       string      `json:"payment_type"`
	Signature_key      string      `json:"signature_key"`
	Status_code        string      `json:"status_code"`
	Status_message     string      `json:"status_message"`
	Transaction_id     string      `json:"transaction_id"`
	Transaction_status string      `json:"transaction_status"`
	Transaction_time   string      `json:"transaction_time"`
	Va_numbers         interface{} `json:"va_numbers"`
	Permata_va_number  string      `json:"permata_va_number"`
}

type VaNumbers struct {
	Bank      string `json:"bank"`
	Va_number string `json:"va_number"`
}

type ConfigAppEnv struct {
	Server_key   string
	Client_key   string
	Midtrans_url string
}
