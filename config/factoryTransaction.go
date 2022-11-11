package config

import (
	"topup-service/helper/payment"
	"topup-service/model"
)

func RepositoryPaymentFactory(dbCon *DatabaseConnection) payment.TransactionList {
	var Repository payment.TransactionList
	if dbCon.Driver == SqlServer {
		Repository = model.NewTransactionModel(dbCon.SqlServer)
	} else {
		panic("Database driver not supported")
	}

	return Repository
}
