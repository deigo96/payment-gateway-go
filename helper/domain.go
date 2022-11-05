package helper

import "time"

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
type GetData struct {
	TotalAllData int
	Result       []Data
}

type Data struct {
	Id                int       `json:"id"`
	User_id           int       `json:"user_id"`
	Reward_Name       string    `json:"reward_name"`
	Description       string    `json:"description"`
	Img_Name          string    `json:"image"`
	Ss_point          int       `json:"point_reward"`
	Request_date      time.Time `json:"request_date"`
	Status_validation int       `json:"status_validation"`
}

type RequestRedeem struct {
	User_id      int `json:"user_id" form:"user_id"`
	Reward_id    int `json:"reward_id" form:"reward_id"`
	Request_date time.Time
}

type TrxPoint struct {
	Id          int
	Trx_type_id int
	Ss_point    int
}

type OwnPoint struct {
	Id             int
	Trx_date       time.Time
	Ss_point_after int
}

type Validate struct {
	Id                int
	User_id           int
	Reward_id         int
	Status_validation int
	Validate_by       int
	Validate_at       time.Time
}

type CurrentPoint struct {
	Id              int
	Trx_date        time.Time
	Ss_point_before int
	Ss_point_trx    int
	Ss_point_after  int
	Exp_date        time.Time
}

type TrxPointList struct {
	Id          int
	Trx_Type_Id int
	Reward_Id   int
	Ss_Point    int
}

type IsUsed struct {
	Is_used int
}

type InsertNewRedeem struct {
	Trx_date         time.Time
	User_id          int
	Is_branch        int
	Trx_point_id     int
	Reference_trx_id int
	Ss_point_before  int
	Ss_point_trx     int
	Ss_point_after   int
}
