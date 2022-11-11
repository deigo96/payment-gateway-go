package modules

import (
	"topup-service/config"
	"topup-service/controller"
	"topup-service/handler"
	csser "topup-service/helper/cart"
	"topup-service/helper/payment"
	"topup-service/router"
)

func RegisterModules(dbCon *config.DatabaseConnection, c *config.AppConfig) router.Controller {
	cart := config.RepositoryPulsaFactory(dbCon)
	cartService := csser.NewCartService(cart)
	transaction := config.RepositoryPaymentFactory(dbCon)
	transactionService := payment.NewTransactionService(transaction)
	jwtService := handler.NewJWTService()
	r := router.Controller{
		Auth: controller.NewAuthController(jwtService, cartService, transactionService),
	}

	return r
}
