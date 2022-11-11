package controller

import (
	"net/http"
	"topup-service/helper"
	"topup-service/helper/cart"

	"github.com/labstack/echo/v4"
)

func (controller *AuthController) GetPackage(c echo.Context) error {
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

	getPackage := controller.cart.GetPackageService()

	resp := helper.BuildSuccessResponse(true, "Succes", getPackage)
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) GetPackageById(c echo.Context) error {
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

	req := new(cart.TopupId)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.NID == 0 {
		resp := helper.BuildErrorResponse("Id topup pack is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	getPackage := controller.cart.GetPackageByIdService(req.NID)

	resp := helper.BuildSuccessResponse(true, "Succes", getPackage)
	return c.JSON(http.StatusOK, resp)
}
