package router

import (
	"topup-service/controller"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Auth *controller.AuthController
}

func Router(e *echo.Echo, controller *Controller) {
	cartRoute := e.Group("api/cart")
	paymentRoute := e.Group("api/payment")
	topUpPack := e.Group("api/package")

	// payment
	paymentRoute.GET("/list-bank", controller.Auth.GetListBank)
	paymentRoute.GET("/payment", controller.Auth.PaymentGateway)
	paymentRoute.GET("/payment/credit-card", controller.Auth.PaymentCreditCartd)
	paymentRoute.GET("/check-transaction", controller.Auth.CheckTransaction)

	//cart
	cartRoute.GET("/get-cart", controller.Auth.GetCart)
	cartRoute.GET("/get-npwp", controller.Auth.GetNpwpUser)
	cartRoute.POST("/add-npwp", controller.Auth.StoreNpwp)
	cartRoute.POST("/add-cart", controller.Auth.StoreCart)
	cartRoute.PUT("/update-cart", controller.Auth.UpdateCart)
	cartRoute.DELETE("/delete-cart", controller.Auth.DeleteCart)

	//package
	topUpPack.GET("/get-package", controller.Auth.GetPackage)
	topUpPack.GET("/get-package-byId", controller.Auth.GetPackageById)
}
