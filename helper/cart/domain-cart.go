package cart

import "time"

type GetData struct {
	TotalAllData int          `json:"total_all_data"`
	Results      []DomainCart `json:"results"`
}

type DomainCart struct {
	NID            int       `json:"id" gorm:"column:nID"`
	Sim            string    `json:"sim"`
	Expired        string    `json:"expired"`
	Bulan_expired  int       `json:"bulan_expired"`
	Next_expired   string    `json:"next_expired"`
	Username       string    `json:"username"`
	Information    string    `json:"information"`
	Date_inserted  time.Time `json:"date_inserted"`
	Harga          int       `json:"harga"`
	Status         int       `json:"status"`
	Order_id       string    `json:"order_id"`
	Domain         string    `json:"domain"`
	Via            string    `json:"via"`
	User_inserted  string    `json:"user_inserted"`
	Sim_status     int       `json:"sim_status"`
	Npwp_id        int       `json:"npwp_id"`
	Top_up_pack_id int       `json:"top_up_pack_id"`
}

type SimCard struct {
	Sim     string
	Expired time.Time
	Status  int
}

type TopupPack struct {
	Price      int
	Topup_days int
}

type Expired struct {
	Expired time.Time
}

type Npwp struct {
	Npwp_id int
}

type NpwpData struct {
	NID        int    `json:"id" gorm:"column:nID"`
	Npwp_no    string `json:"npwp_no"`
	Npwp_name  string `json:"npwp_name"`
	Npwp_addr  string `json:"npwp_addr"`
	Npwp_wa    string `json:"npwp_wa"`
	Npwp_email string `json:"npwp_email"`
	Created_at string `json:"created_at"`
}

type Package struct {
	TotalAllData int         `json:"total_all_data"`
	Results      []TopupData `json:"results"`
}

type TopupData struct {
	NID        int    `json:"id" gorm:"column:nID"`
	Pack_name  string `json:"pack_name"`
	Price      string `json:"price"`
	Topup_days string `json:"topup_days"`
	Is_default int    `json:"is_default"`
}

type TopupId struct {
	NID int `json:"id" gorm:"column:nID"`
}
