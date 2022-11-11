package cart

import "time"

type Data struct {
	User           UserData `json:"user"`
	Vehicle        Vehicle  `json:"vehicle"`
	Top_up_pack_id int      `json:"top_up_pack_id"`
	Is_mobile      int      `json:"is_mobile"`
}

type UserData struct {
	Username  string `json:"username"`
	Domain    string `json:"domain"`
	Full_name string `json:"fullname"`
}

type RequestCart struct {
	Sim            string    `json:"sim"`
	Expired        string    `json:"expired"`
	Bulan_expired  int       `json:"bulan_expired"`
	Next_expired   string    `json:"next_expired"`
	Plate          string    `json:"plate"`
	Username       string    `json:"username"`
	Information    string    `json:"information"`
	Date_inserted  time.Time `json:"date_inserted"`
	Harga          int       `json:"harga"`
	Status         int       `json:"status"`
	Order_id       string    `json:"order_id"`
	Domain         string    `json:"domain"`
	Via            string    `json:"via"`
	User_inserted  string    `json:"user_inserted"`
	Npwp_id        int       `json:"npwp_id"`
	Top_up_pack_id int       `json:"top_up_pack_id"`
	Is_web         int       `json:"is_web"`
}

type Vehicle struct {
	Plate string `json:"plate"`
	Sim   string `json:"sim"`
}

type RequestUpdate struct {
	NID            int `json:"cart_id" gorm:"column:nId"`
	Is_mobile      int `json:"is_mobile"`
	Top_up_pack_id int `json:"top_up_pack_id"`
}

type RequestDelete struct {
	NID int `json:"cart_id"`
}

type UpdateCartForPayment struct {
	Status   int    `json:"status"`
	Order_id string `json:"order_id"`
}

type RequestNpwp struct {
	ID         uint      `gorm:"column:nID"`
	Username   string    `json:"username"`
	Npwp_no    string    `json:"npwp_no"`
	Npwp_name  string    `json:"npwp_name"`
	Npwp_addr  string    `json:"npwp_addr"`
	Npwp_wa    string    `json:"npwp_wa"`
	Npwp_email string    `json:"npwp_email"`
	Created_at time.Time `json:"created_at"`
}
