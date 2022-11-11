package controller

import (
	"net/http"
	"topup-service/helper"

	"github.com/labstack/echo/v4"
)

func (controller *AuthController) RecurringXendit(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	if header == "Bearer " || header == "" {
		resp := helper.BuildErrorResponse("No Authorization, No Token", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	validate := controller.jwtService.ValidateToken(header)
	if validate != true {
		resp := helper.BuildErrorResponse("Token expired or invalid", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}
	resp := helper.BuildErrorResponse("coreApiRes.StatusMessage, a", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}
