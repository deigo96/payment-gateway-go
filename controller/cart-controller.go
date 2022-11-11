package controller

import (
	"net/http"
	"strconv"
	"time"
	"topup-service/helper"
	"topup-service/helper/cart"

	"github.com/labstack/echo/v4"
)

func (controller *AuthController) GetCart(c echo.Context) (err error) {
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

	username := c.QueryParam("username")
	domain := c.QueryParam("domain")

	cartData, err := controller.cart.GetCartService(username, domain)
	if err != nil {
		resp := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	resp := helper.BuildSuccessResponse(true, "successful get cart", cartData)
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) GetNpwpUser(c echo.Context) (err error) {
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

	username := c.QueryParam("username")
	domain := c.QueryParam("domain")

	cartData, err := controller.cart.GetNpwpService(username, domain)
	if err != nil {
		resp := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusUnauthorized, resp)
	}

	resp := helper.BuildSuccessResponse(true, "successful get cart", cartData)
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) StoreNpwp(c echo.Context) (err error) {
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

	req := new(cart.RequestNpwp)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	if req.Username == "" {
		resp := helper.BuildErrorResponse("Username is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Npwp_no == "" {
		resp := helper.BuildErrorResponse("Npwp No is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Npwp_addr == "" {
		resp := helper.BuildErrorResponse("Npwp address is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Npwp_email == "" {
		resp := helper.BuildErrorResponse("Npwp email is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Npwp_name == "" {
		resp := helper.BuildErrorResponse("Npwp name is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	if req.Npwp_wa == "" {
		resp := helper.BuildErrorResponse("Npwp wa is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	store := cart.RequestNpwp{
		Username:   req.Username,
		Npwp_no:    req.Npwp_no,
		Npwp_name:  req.Npwp_name,
		Npwp_addr:  req.Npwp_addr,
		Npwp_email: req.Npwp_email,
		Npwp_wa:    req.Npwp_wa,
		Created_at: time.Now(),
	}

	err = controller.cart.StoreNpwpService(store)
	if err != nil {
		resp := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusOK, resp)
	}

	resp := helper.BuildSuccessResponse(true, "Successful added npwp", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) StoreCart(c echo.Context) (err error) {
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

	req := new(cart.Data)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	// if req.Sim == "" {
	// 	response := helper.BuildErrorResponse("Sim is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Expired == "" {
	// 	response := helper.BuildErrorResponse("Expired is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Expired == "" {
	// 	response := helper.BuildErrorResponse("Expired is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Username == "" {
	// 	response := helper.BuildErrorResponse("Username is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Plate == "" {
	// 	response := helper.BuildErrorResponse("Plate is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	info := req.Vehicle.Plate + "|" + req.User.Full_name + "|" + req.Vehicle.Sim

	via := "web"
	if req.Is_mobile == 1 {
		via = "mobile"
	}

	unix := time.Now().UnixMilli()

	store := cart.RequestCart{
		Sim:            req.Vehicle.Sim,
		Username:       req.User.Username,
		Information:    info,
		Date_inserted:  time.Now(),
		Status:         0,
		Order_id:       "INV-001" + strconv.Itoa(int(unix)),
		Domain:         req.User.Domain,
		Via:            via,
		Top_up_pack_id: req.Top_up_pack_id,
	}
	err = controller.cart.StoreCartService(store)

	if err != nil {
		response := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusOK, response)
	}

	resp := helper.BuildSuccessResponse(true, "successfully added to cart", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) UpdateCart(c echo.Context) (err error) {
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

	req := new(cart.RequestUpdate)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	// if req.Bulan_expired == 0 {
	// 	response := helper.BuildErrorResponse("Bulan expired is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Next_expired == "" {
	// 	response := helper.BuildErrorResponse("Next expired is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Harga == 0 {
	// 	response := helper.BuildErrorResponse("Harga is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Top_up_pack_id == 0 {
	// 	response := helper.BuildErrorResponse("Top up pack id is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	// if req.Order_id == "" {
	// 	response := helper.BuildErrorResponse("Order_id is required", helper.EmptyObj{})
	// 	return c.JSON(http.StatusOK, response)
	// }

	update := cart.RequestUpdate{
		Is_mobile:      req.Is_mobile,
		Top_up_pack_id: req.Top_up_pack_id,
	}

	err = controller.cart.UpdateCartService(update, req.NID)
	if err != nil {
		response := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusOK, response)
	}

	resp := helper.BuildSuccessResponse(true, "Successfully updated cart", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}

func (controller *AuthController) DeleteCart(c echo.Context) (err error) {
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

	req := new(cart.RequestDelete)
	if err = c.Bind(req); err != nil {
		response := helper.BuildErrorResponse("Invalid request body", helper.EmptyObj{})
		return c.JSON(http.StatusBadRequest, response)
	}

	if req.NID == 0 {
		response := helper.BuildErrorResponse("Cart id is required", helper.EmptyObj{})
		return c.JSON(http.StatusOK, response)
	}

	err = controller.cart.DeleteCartService(req.NID)
	if err != nil {
		response := helper.BuildErrorResponse(err.Error(), helper.EmptyObj{})
		return c.JSON(http.StatusOK, response)
	}

	resp := helper.BuildSuccessResponse(true, "Successfully deleted cart", helper.EmptyObj{})
	return c.JSON(http.StatusOK, resp)
}
