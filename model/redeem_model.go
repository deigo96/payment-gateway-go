package model

import (
	"errors"
	"request-redeem/helper"
	"time"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

type RewardRepo struct {
	db *gorm.DB
}

func NewStoreReward(db *gorm.DB) helper.RewardList {
	return &RewardRepo{
		db: db,
	}
}

func (r *RewardRepo) GetRequestRedeem() (d []helper.Data) {
	p := r.db.Table("request_redeem rd").Select("rd.id, rd.user_id, rm.reward_name, rm.description, rm.img_name, tp.ss_point, rd.request_date, rd.status_validation").Joins("left join reward_master rm on rd.reward_id = rm.id").Joins("left join transaction_point tp on rd.reward_id = tp.reward_id").Find(&d)
	if p.Error != nil {
		return d
	}

	return d
}

func (r *RewardRepo) GetRequestRedeemById(id int) (d []helper.Data) {
	p := r.db.Table("request_redeem rd").Select("rd.id, rd.user_id, rm.reward_name, rm.description, rm.img_name, tp.ss_point, rd.request_date, rd.status_validation").Joins("left join reward_master rm on rd.reward_id = rm.id").Joins("left join transaction_point tp on rd.reward_id = tp.reward_id").Where("rd.user_id = ?", id).Find(&d)
	if p.Error != nil {
		return d
	}

	return d
}

func (r *RewardRepo) StoreRequestRedeem(s helper.RequestRedeem) (err error) {
	var point helper.TrxPoint
	p := r.db.Table("transaction_point t").Select("t.id, t.trx_type_id, t.ss_point").Where("reward_id = ?", s.Reward_id).Find(&point)
	if p.Error != nil {
		return p.Error
	}

	var own helper.OwnPoint
	o := r.db.Table("referral_transaction r").Select("r.id, r.trx_date, r.ss_point_after").Where("user_id = ? AND is_expired = ? AND is_used = ?", s.User_id, 0, 0).Order("r.id DESC LIMIT 1").Find(&own)
	if o.Error != nil {
		return o.Error
	}

	if point.Ss_point > own.Ss_point_after {
		return errors.New("Not enough point")
	}

	c := r.db.Table("request_redeem").Find(&s)
	if c.RowsAffected > 0 {
		return errors.New("Request redeem already send")
	}

	res := r.db.Table("request_redeem").Create(&s)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *RewardRepo) ValidateRedeem(s helper.Validate) (err error) {
	var check helper.Validate
	firstMonth := now.BeginningOfMonth()
	// lastMonth := now.EndOfMonth()

	c := r.db.Table("request_redeem").Where("id = ?", s.Id).Find(&check)
	if c.RowsAffected == 0 {
		return errors.New("Id request redeem not found")
	}
	if check.Status_validation != 0 {
		return errors.New("You already process the request")
	}

	if s.Status_validation == 1 {
		var cp helper.CurrentPoint
		var dataTrx helper.TrxPointList

		data := r.db.Table("transaction_point").Select("id, trx_type_id, reward_id, ss_point, own_pct, branch_pct, referral_pct").Where("reward_id = ?", check.Reward_id).Find(&dataTrx)
		if data.Error != nil {
			return data.Error
		}

		currentPoint := r.db.Table("referral_transaction rt").Select("rt.id, rt.trx_date,rt.ss_point_before, rt.ss_point_after, rt.ss_point_trx, rt.exp_date").Where("rt.user_id = ? AND exp_date >= ?", check.User_id, firstMonth).Order("rt.id DESC LIMIT 1").Find(&cp)
		if currentPoint.Error != nil {
			return currentPoint.Error
		}

		if dataTrx.Ss_Point > cp.Ss_point_after {
			return errors.New("Not enough point")
		}

		pointAfter := cp.Ss_point_after - dataTrx.Ss_Point

		insert := helper.InsertNewRedeem{
			Trx_date:         time.Now(),
			User_id:          check.User_id,
			Trx_point_id:     dataTrx.Id,
			Reference_trx_id: check.Reward_id,
			Ss_point_before:  cp.Ss_point_after,
			Ss_point_trx:     -(dataTrx.Ss_Point),
			Ss_point_after:   pointAfter,
		}

		insertRedeem := r.db.Table("referral_transaction").Create(&insert)
		if insertRedeem.Error != nil {
			return insertRedeem.Error
		}

	}

	u := r.db.Table("request_redeem").Where("id = ?", s.Id).Updates(&s)
	if u.Error != nil {
		return u.Error
	}

	return nil
}
