package config

import (
	"topup-service/helper/cart"
	"topup-service/model"
)

func RepositoryPulsaFactory(dbCon *DatabaseConnection) cart.CartList {
	var Repository cart.CartList
	if dbCon.Driver == SqlServer {
		Repository = model.NewCartModel(dbCon.SqlServer)
	} else {
		panic("Database driver not supported")
	}

	return Repository
}
