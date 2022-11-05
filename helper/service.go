package helper

type RewardList interface {
	GetRequestRedeem() (d []Data)
	GetRequestRedeemById(id int) (d []Data)
	StoreRequestRedeem(s RequestRedeem) (err error)
	ValidateRedeem(s Validate) (err error)
}

type RewardService interface {
	GetRequestRedeemService() (d []Data)
	GetRequestRedeemByIdService(id int) (d []Data)
	StoreRequestRedeemService(s RequestRedeem) (err error)
	ValidateRedeemService(s Validate) (err error)
}

type rewardService struct {
	reward RewardList
}

func NewRewardService(reward RewardList) RewardService {
	return &rewardService{
		reward: reward,
	}
}

func (r *rewardService) GetRequestRedeemService() (d []Data) {
	res := r.reward.GetRequestRedeem()

	return res
}

func (r *rewardService) GetRequestRedeemByIdService(id int) (d []Data) {
	res := r.reward.GetRequestRedeemById(id)

	return res
}

func (r *rewardService) StoreRequestRedeemService(s RequestRedeem) (err error) {
	err = r.reward.StoreRequestRedeem(s)
	if err != nil {
		return err
	}

	return nil
}

func (r *rewardService) ValidateRedeemService(s Validate) (err error) {
	err = r.reward.ValidateRedeem(s)
	if err != nil {
		return err
	}
	return nil
}
