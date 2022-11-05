package router

import (
	"request-redeem/controller"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Auth *controller.AuthController
}

func Routes(e *echo.Echo, controller *Controller) {
	e.GET("/requestRedeem", controller.Auth.GetRequestRedeem)
	e.POST("/requestRedeem", controller.Auth.RequestRedeem)
	e.PUT("/requestRedeem", controller.Auth.ValidateRequest)
}
