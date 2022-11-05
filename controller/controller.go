package controller

import (
	"net/http"
	"request-redeem/handler"
	"request-redeem/helper"
	"request-redeem/helper/api"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	Service    helper.RewardService
	jwtService handler.JwtService
}

func NewAuthController(Service helper.RewardService, jwtService handler.JwtService) *AuthController {
	return &AuthController{
		Service:    Service,
		jwtService: jwtService,
	}
}

func (controller *AuthController) GetRequestRedeem(c echo.Context) (err error) {
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

	jwtPrivilege := controller.jwtService.ValidatePrivilege(header)
	md, _ := jwtPrivilege.Data.(map[string]interface{})
	privilege := md["privilege"].(float64)

	if int(privilege) == 1 || int(privilege) == 2 {
		redeem := controller.Service.GetRequestRedeemService()
		if len(redeem) < 1 {
			resp := helper.BuildErrorResponse("No request redeem", helper.EmptyObj{})
			return c.JSON(http.StatusBadRequest, resp)
		}
		resp := helper.BuildSuccessResponse(true, "Success get request redeem", redeem)
		return c.JSON(http.StatusOK, resp)
	}

	seen, err := api.GetUserSeen(header)
	if err != nil {
		resp := helper.BuildErrorResponse("User not found", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	redeem := controller.Service.GetRequestRedeemByIdService(seen.Data.Id)
	if len(redeem) < 1 {
		resp := helper.BuildErrorResponse("No request redeem", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp := helper.BuildSuccessResponse(true, "Success get request redeem", redeem)
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) RequestRedeem(c echo.Context) (err error) {
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

	seen, err := api.GetUserSeen(header)
	if err != nil {
		resp := helper.BuildErrorResponse("User not found", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	req := new(helper.RequestRedeem)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	if req.Reward_id == 0 {
		resp := helper.BuildErrorResponse("Reward id is required", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}

	store := helper.RequestRedeem{
		Request_date: time.Now(),
		User_id:      seen.Data.Id,
		Reward_id:    req.Reward_id,
	}

	err = controller.Service.StoreRequestRedeemService(store)
	if err != nil {
		resp := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp := helper.BuildSuccessResponse(true, "Success request redeem reward", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) ValidateRequest(c echo.Context) (err error) {
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

	jwtPrivilege := controller.jwtService.ValidatePrivilege(header)
	md, _ := jwtPrivilege.Data.(map[string]interface{})
	privilege := md["privilege"].(float64)

	seen, err := api.GetUserSeen(header)
	if err != nil {
		resp := helper.BuildErrorResponse("User not found", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	if int(privilege) != 1 && int(privilege) != 2 {
		resp := helper.BuildErrorResponse("Unauthorized", helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	req := new(helper.Validate)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	if req.Id == 0 {
		resp := helper.BuildErrorResponse("Id request is required", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}

	if req.Status_validation == 0 {
		resp := helper.BuildErrorResponse("Status validation is required", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}

	store := helper.Validate{
		Id:                req.Id,
		Status_validation: req.Status_validation,
		Validate_by:       seen.Data.Id,
		Validate_at:       time.Now(),
	}

	err = controller.Service.ValidateRedeemService(store)
	if err != nil {
		resp := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp := helper.BuildSuccessResponse(true, "Success validate request redeem", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}
