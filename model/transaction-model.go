package model

import (
	"topup-service/helper"
	"topup-service/helper/payment"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionModel(db *gorm.DB) payment.TransactionList {
	return &TransactionRepo{
		db: db,
	}
}

func (t *TransactionRepo) StoreTransactionList(s helper.RequestStoreTransaction) error {
	res := t.db.Table("pulsa_dev.dbo.store_transaction").Create(&s)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (t *TransactionRepo) CheckPreviousTrxList(username string) (s []helper.RequestStatus, err error) {
	res := t.db.Table("pulsa_dev.dbo.store_transaction").Where("username = ? AND status = 'pending'", username).Find(&s)
	if res.Error != nil {
		return s, nil
	}

	status := helper.UpdateStatusTrx{
		Status: "cancel",
	}

	for _, val := range s {
		_ = t.db.Table("pulsa_dev.dbo.store_transaction").Where("order_id = ?", val.Order_id).Updates(&status)

	}

	return s, nil
}

func (t *TransactionRepo) UpdateStoreTransactionList(s helper.RequestStoreTransaction, orderId string, status int) error {
	data := map[string]interface{}{
		"status":           s.Status,
		"mid_json_payment": s.Mid_json_payment,
	}
	res := t.db.Table("pulsa_dev.dbo.store_transaction").Where("order_id = ?", orderId).Updates(&data)
	if res.Error != nil {
		return res.Error
	}

	dataCart := map[string]interface{}{
		"status": status,
	}
	_ = t.db.Table("pulsa_dev.dbo.cart").Where("order_id = ?", orderId).Updates(&dataCart)

	return nil
}
