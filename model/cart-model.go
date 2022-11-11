package model

import (
	"errors"
	"fmt"
	"topup-service/helper"
	"topup-service/helper/cart"

	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartModel(db *gorm.DB) cart.CartList {
	return &CartRepo{
		db: db,
	}
}

func (r *CartRepo) GetCartList(username string, domain string) (s []cart.DomainCart, err error) {

	res := r.db.Table("pulsa_dev.dbo.cart c").Where("c.username = ? AND c.status = 0 AND c.domain = ?", username, domain).Order("date_inserted DESC").Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	return s, nil
}

func (r *CartRepo) GetNpwpList(username string, domain string) (s cart.NpwpData, err error) {
	var a cart.Npwp

	res := r.db.Table("pulsa_dev.dbo.npwp n").Where("n.nID IN (?)", r.db.Table("pulsa_dev.dbo.cart c").Select("c.npwp_id").Where("c.username = ? AND c.domain = ? AND npwp_id is not null", username, domain).Order("c.nID DESC").Limit(1).Find(&a)).Find(&s)
	if res.Error != nil {
		return s, res.Error
	}

	return s, nil
}

func (r *CartRepo) StoreNpwpList(s cart.RequestNpwp) error {

	res := r.db.Table("pulsa_dev.dbo.npwp").Create(&s)
	if res.Error != nil {
		return res.Error
	}
	npwp := cart.Npwp{
		Npwp_id: int(s.ID),
	}

	_ = r.db.Table("pulsa_dev.dbo.cart").Where("nID IN (?)", r.db.Table("pulsa_dev.dbo.cart").Select("nID").Where("username = ?", s.Username).Order("date_inserted DESC").Limit(1)).Updates(&npwp)

	return nil
}

func (r *CartRepo) StoreCartList(s cart.RequestCart) (err error) {
	var c cart.SimCard
	var a cart.TopupPack
	_ = r.db.Table("pulsa_dev.dbo.SimCard").Where("sim = ?", s.Sim).Find(&c)
	_ = r.db.Table("pulsa_dev.dbo.topup_packs").Where("nID = ?", s.Top_up_pack_id).Find(&a)
	bulanExpired := a.Topup_days / 30
	nextExpired := c.Expired.AddDate(0, 0, a.Topup_days)

	req := map[string]interface{}{
		"sim":            s.Sim,
		"expired":        c.Expired,
		"bulan_expired":  bulanExpired,
		"next_expired":   nextExpired,
		"username":       s.Username,
		"information":    s.Information,
		"date_inserted":  s.Date_inserted,
		"harga":          a.Price,
		"status":         s.Status,
		"order_id":       s.Order_id,
		"domain":         s.Domain,
		"via":            s.Via,
		"top_up_pack_id": s.Top_up_pack_id,
		"sim_status":     c.Status,
		"privilege_id":   0,
	}

	res := r.db.Table("pulsa_dev.dbo.cart").Create(&req)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *CartRepo) UpdateCartList(s cart.RequestUpdate, id int) (err error) {
	var c cart.Expired
	var a cart.TopupPack
	check := r.db.Table("pulsa_dev.dbo.cart").Where("nID = ?", id).Find(&c)
	if check.RowsAffected == 0 {
		return errors.New("Cart id not found")
	}
	_ = r.db.Table("pulsa_dev.dbo.topup_packs").Where("nID = ?", s.Top_up_pack_id).Find(&a)
	bulanExpired := a.Topup_days / 30
	nextExpired := c.Expired.AddDate(0, 0, a.Topup_days)

	via := "web"
	if s.Is_mobile == 1 {
		via = "mobile"
	}

	req := map[string]interface{}{
		"bulan_expired":  bulanExpired,
		"next_expired":   nextExpired,
		"harga":          a.Price,
		"top_up_pack_id": s.Top_up_pack_id,
		"via":            via,
	}

	res := r.db.Table("pulsa_dev.dbo.cart").Where("nID = ? AND status = 0", id).Updates(&req)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("No data found")
	}

	return nil

}

func (r *CartRepo) DeleteCartList(id int) (err error) {
	var d cart.RequestDelete

	chekcOrderId := r.db.Table("pulsa_dev.dbo.cart").Where("status = 0 AND nID = ?", id).Find(&d)
	if chekcOrderId.RowsAffected == 0 {
		return errors.New("Cart id not found")
	}

	delete := r.db.Table("pulsa_dev.dbo.cart").Where("status = 0 AND nID = ?", id).Delete(&d)
	if delete.Error != nil {
		return delete.Error
	}

	return nil
}

func (r *CartRepo) UpdateCartPaymentList(s helper.RequestUpdateCart) error {
	data := map[string]interface{}{
		"status":   s.Status,
		"order_id": s.Order_id,
	}

	_ = r.db.Table("pulsa_dev.dbo.cart").Where("status = 0 AND username = ? AND domain = ?", s.Username, s.Domain).Updates(&data)

	return nil
}

func (r *CartRepo) UpdateStatusCartPaymentList(order_id string) error {
	data := map[string]interface{}{
		"status": 3,
	}
	res := r.db.Table("pulsa_dev.dbo.cart").Where("order_id = ?", order_id).Updates(&data)
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	return nil
}

func (r *CartRepo) GetPackageList() (s []cart.TopupData) {
	if err := r.db.Table("pulsa_dev.dbo.topup_packs").Order("nID ASC").Find(&s).Error; err != nil {
		return s
	}
	return s
}

func (r *CartRepo) GetPackageByIdList(id int) (s cart.TopupData) {
	if err := r.db.Table("pulsa_dev.dbo.topup_packs").Where("nID = ?", id).Find(&s).Error; err != nil {
		return s
	}
	return s
}
