package goswyftx

import "strconv"

type HistoryService struct {
	service
	assetId int
}

type CurrencyHistory struct {
	ID        int        `json:"id,omitempty"`
	Time      SwyftxTime `json:"time,omitempty"`
	Quantity  string     `json:"quantity,omitempty"`
	AddressID int        `json:"address_id,omitempty"`
	Status    string     `json:"status,omitempty"`
}

type TransactionHistory struct {
	Asset      int        `json:"asset,omitempty"`
	Amount     float32    `json:"amount,omitempty"`
	Updated    SwyftxTime `json:"updated,omitempty"`
	ActionType string     `json:"actionType,omitempty"`
	Status     string     `json:"status,omitempty"`
}

// History will return a history service that holds methods which can get history events for an asset
func (c *Client) History(asset int) *HistoryService {
	return &HistoryService{service{c}, asset}
}

// Withdraw events for an asset
func (hs *HistoryService) Withdraw() (*CurrencyHistory, error) {
	return hs.currency("withdraw")
}

// Deposit events for an asset
func (hs *HistoryService) Deposit() (*CurrencyHistory, error) {
	return hs.currency("deposit")
}

func (hs *HistoryService) currency(actionType string) (*CurrencyHistory, error) {
	var histCurrency CurrencyHistory
	if err := hs.client.Get(buildString("history/", actionType, "/", strconv.Itoa(hs.assetId)),
		&histCurrency); err != nil {
		return nil, err
	}

	return &histCurrency, nil
}

// All trades, withdrawals and deposits events for an asset
func (hs *HistoryService) All(actionType string) ([]*TransactionHistory, error) {
	var transHist []*TransactionHistory
	if err := hs.client.Get(buildString("history/", actionType, "/", strconv.Itoa(hs.assetId)),
		&transHist); err != nil {
		return nil, err
	}

	return transHist, nil
}
