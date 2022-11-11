package cart

import "topup-service/helper"

type CartList interface {
	GetCartList(username string, domain string) (s []DomainCart, err error)
	GetNpwpList(username string, domain string) (s NpwpData, err error)
	StoreCartList(s RequestCart) error
	UpdateCartList(s RequestUpdate, id int) error
	DeleteCartList(id int) error
	StoreNpwpList(RequestNpwp) error
	UpdateCartPaymentList(helper.RequestUpdateCart) error
	UpdateStatusCartPaymentList(order_id string) error
	GetPackageList() []TopupData
	GetPackageByIdList(id int) TopupData
}

type CartService interface {
	GetCartService(username string, domain string) (s GetData, err error)
	GetNpwpService(username string, domain string) (s NpwpData, err error)
	StoreCartService(RequestCart) error
	UpdateCartService(s RequestUpdate, id int) error
	DeleteCartService(id int) error
	StoreNpwpService(RequestNpwp) error
	UpdateCartPaymentService(helper.RequestUpdateCart) error
	UpdateStatusCartPaymentService(order_id string) error
	GetPackageService() Package
	GetPackageByIdService(id int) TopupData
}

type cartService struct {
	cart CartList
}

func NewCartService(cart CartList) CartService {
	return &cartService{
		cart: cart,
	}
}

func (c *cartService) GetCartService(username string, domain string) (s GetData, err error) {
	res, err := c.cart.GetCartList(username, domain)
	if err != nil {
		return s, err
	}
	s.TotalAllData = len(res)
	s.Results = append(s.Results, res...)

	return s, nil
}

func (c *cartService) GetNpwpService(username string, domain string) (s NpwpData, err error) {
	res, err := c.cart.GetNpwpList(username, domain)
	if err != nil {
		return s, err
	}

	return res, nil
}

func (c *cartService) StoreNpwpService(s RequestNpwp) (err error) {
	err = c.cart.StoreNpwpList(s)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartService) StoreCartService(s RequestCart) (err error) {
	err = c.cart.StoreCartList(s)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartService) UpdateCartService(s RequestUpdate, id int) (err error) {
	err = c.cart.UpdateCartList(s, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartService) DeleteCartService(id int) (err error) {
	err = c.cart.DeleteCartList(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartService) UpdateCartPaymentService(s helper.RequestUpdateCart) error {
	_ = c.cart.UpdateCartPaymentList(s)
	return nil
}

func (c *cartService) UpdateStatusCartPaymentService(order_id string) error {
	_ = c.cart.UpdateStatusCartPaymentList(order_id)
	return nil
}

func (c *cartService) GetPackageService() (s Package) {
	res := c.cart.GetPackageList()
	s.Results = res
	s.TotalAllData = len(res)

	return s
}

func (c *cartService) GetPackageByIdService(id int) (s TopupData) {
	res := c.cart.GetPackageByIdList(id)

	return res
}
