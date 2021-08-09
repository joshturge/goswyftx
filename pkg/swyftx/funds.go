package goswyftx

import "strconv"

type FundsService struct {
	service
	id int
}

// Funds will return a funds service that holds a method which can withdraw funds from an address
func (c *Client) Funds(addressId int) *FundsService {
	return &FundsService{service{c}, addressId}
}

// Withdraw funds from an account into a specified asset
func (fs *FundsService) Withdraw(asset int, amount float32) error {
	var body struct {
		Quantity  float32 `json:"quantity"`
		AddressID int     `json:"address_id"`
	}
	body.AddressID = fs.id
	body.Quantity = amount

	if err := fs.client.Post(buildString("funds/withdraw/", strconv.Itoa(asset)), &body, nil); err != nil {
		return err
	}

	return nil
}
