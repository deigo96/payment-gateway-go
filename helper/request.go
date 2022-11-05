package helper

type StoreRewardRequest struct {
	Reward_Name        string `json:"rewardName" form:"rewardName"`
	Reward_Image       string `json:"rewardImage" form:"rewardImage"`
	Reward_Description string `json:"rewardDescription" form:"rewardDescription"`
	Reward_category_id int    `json:"idCategory" form:"rewardCategoryId"`
}

type Idreward struct {
	Id_Reward int `json:"idReward" form:"idReward"`
}
