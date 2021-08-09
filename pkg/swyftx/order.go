package swyftx

import "strconv"

type OrderService service

type OrderExchangeRate struct {
	Mid   string `json:"mid,omitempty"`
	Price string `json:"price,omitempty"`
}

type OrderPlace struct {
	Primary       string  `json:"primary,omitempty"`
	Secondary     string  `json:"secondary,omitempty"`
	Quantity      float32 `json:"quantity,omitempty"`
	AssetQuantity string  `json:"assetQuantity,omitempty"`
	OrderType     string  `json:"orderType,omitempty"`
	Trigger       int     `json:"trigger,omitempty"`
}

type Order struct {
	Type           string     `json:"order_type,omitempty"`
	PrimaryAsset   string     `json:"primary_asset,omitempty"`
	SecondaryAsset string     `json:"secondary_asset,omitempty"`
	QuantityAsset  string     `json:"quantity_asset,omitempty"`
	Quantity       int        `json:"quantity,omitempty"`
	Trigger        int        `json:"trigger,omitempty"`
	Status         string     `json:"status,omitempty"`
	Amount         int        `json:"amount,omitempty"`
	Total          float32    `json:"total,omitempty"`
	Price          int        `json:"price,omitempty"`
	CreateTime     SwyftxTime `json:"created_time,omitempty"`
	ID             int        `json:"id,omitempty"`
}

// Order will return a order service that can interact with swyftx api
func (c *Client) Order() *OrderService {
	return (*OrderService)(&service{c})
}

// PairExchangeRate will show the exchange rate for a crypto pair
func (os *OrderService) PairExchangeRate(buy, sell string, amount int,
	limit string) (*OrderExchangeRate, error) {
	var (
		body struct {
			Buy    string `json:"buy,omitempty"`
			Sell   string `json:"sell,omitempty"`
			Amount int    `json:"amount,omitempty"`
			Limit  string `json:"limit,omitempty"`
		}
		exchRate OrderExchangeRate
	)
	body.Buy = buy
	body.Sell = sell
	body.Amount = amount
	body.Limit = limit

	if err := os.client.Post("orders/rate/", &body, &exchRate); err != nil {
		return nil, err
	}

	return &exchRate, nil
}

// Place will create an order from an OrderPlace, returns the order id
func (os *OrderService) Place(order *OrderPlace) (int, error) {
	var orderID struct {
		OrderID int `json:"orderId"`
	}
	if err := os.client.Post("orders/", order, &orderID); err != nil {
		return 0, err
	}

	return orderID.OrderID, nil
}

// Cancel will cancel an order
func (os *OrderService) Cancel(orderID int) error {
	if err := os.client.Delete(buildString("orders/", strconv.Itoa(orderID))); err != nil {
		return err
	}

	return nil
}

// List all orders for an asset
func (os *OrderService) List(asset string) ([]*Order, error) {
	var orders []*Order
	if err := os.client.Get(buildString("orders/", asset), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
