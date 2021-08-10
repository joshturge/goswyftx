package swyftx

import "strconv"

// MarketService contains method that interact with the market swyftx API
// endpoints
type MarketService service

// MarketRate contains information about the current price of an asset against
// an asset
type MarketRate struct {
	DailyPriceChange string `json:"dailyPriceChange,omitempty"`
	MidPrice         string `json:"midPrice,omitempty"`
}

// MarketAsset contains information about an asset provided by swyftx
type MarketAsset struct {
	ID                    int     `json:"id,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Code                  string  `json:"code,omitempty"`
	MinimumOrder          string  `json:"minimum_order,omitempty"`
	PriceScale            int     `json:"price_scale,omitempty"`
	DepositEnabled        bool    `json:"deposit_enabled,omitempty"`
	WithdrawEnabled       bool    `json:"withdraw_enabled,omitempty"`
	MinConfirmations      int     `json:"min_confirmations,omitempty"`
	MinWithdrawal         int     `json:"min_withdrawal,omitempty"`
	MinimumOrderIncrement float32 `json:"minimum_order_increment,omitempty"`
	MiningFee             int     `json:"mining_fee,omitempty"`
	Primary               bool    `json:"primary,omitempty"`
	Secondary             bool    `json:"secondary,omitempty"`
}

// MarketBasicInfo contains basic information about an asset
type MarketBasicInfo struct {
	Name      string  `json:"name,omitempty"`
	AltName   string  `json:"altName,omitempty"`
	Code      string  `json:"code,omitempty"`
	ID        int     `json:"id,omitempty"`
	Rank      int     `json:"rank,omitempty"`
	Buy       string  `json:"buy,omitempty"`
	Sell      string  `json:"sell,omitempty"`
	Spread    string  `json:"spread,omitempty"`
	Volume24H float32 `json:"volume24H,omitempty"`
	MarketCap float64 `json:"marketCap,omitempty"`
}

// MarketDetailedInfo contains detailed information about an asset
type MarketDetailedInfo struct {
	Name        string `json:"name,omitempty"`
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Mineable    int    `json:"mineable,omitempty"`
	Spread      string `json:"spread,omitempty"`
	Rank        int    `json:"rank,omitempty"`
	RankSuffix  string `json:"rankSuffix,omitempty"`
	Volume      struct {
		H24       float32 `json:"24H,omitempty"`
		W1        float32 `json:"1W,omitempty"`
		M1        float32 `json:"1M,omitempty"`
		MarketCap float64 `json:"marketCap,omitempty"`
	} `json:"volume,omitempty"`
	URL struct {
		Website  string `json:"website,omitempty"`
		Twitter  string `json:"twitter,omitempty"`
		Reddit   string `json:"reddit,omitempty"`
		TechDoc  string `json:"techDoc,omitempty"`
		Explorer string `json:"explorer,omitempty"`
	} `json:"urls,omitempty"`
	Supply struct {
		Circulating int `json:"circulating,omitempty"`
		Total       int `json:"total,omitempty"`
		Max         int `json:"max,omitempty"`
	} `json:"supply,omitempty"`
}

// Market will return a market service which can interact with swyftx API market
// endpoints
func (c *Client) Market() *MarketService {
	return (*MarketService)(&service{c})
}

// LiveRates will get live rates from swyftx for a given asset
func (ms *MarketService) LiveRates(asset int) (*MarketRate, error) {
	var marketRate struct {
		// TODO: fix this, their is more than one asset
		MarketRate MarketRate `json:"1"`
	}
	if err := ms.client.Get(buildString("live-rates/", strconv.Itoa(asset)), &marketRate); err != nil {
		return nil, err
	}

	return &marketRate.MarketRate, nil
}

// Assets will get market information on all swyftx assets
func (ms *MarketService) Assets() ([]*MarketAsset, error) {
	var marketAssets []*MarketAsset
	if err := ms.client.Get("markets/assets/", &marketAssets); err != nil {
		return nil, err
	}

	return marketAssets, nil
}

// BasicInfo will get basic information for a given an asset code
func (ms *MarketService) BasicInfo(assetCode string) (*MarketBasicInfo, error) {
	var marketBasic MarketBasicInfo
	if err := ms.client.Get(buildString("markets/info/basic/", assetCode), &marketBasic); err != nil {
		return nil, err
	}

	return &marketBasic, nil
}

// DetailedInfo will get detailed information on an asset given an asset code
func (ms *MarketService) DetailedInfo(assetCode string) ([]*MarketDetailedInfo, error) {
	var detailedInfo []*MarketDetailedInfo
	if err := ms.client.Get(buildString("markets/info/details/", assetCode), &detailedInfo); err != nil {
		return nil, err
	}

	return detailedInfo, nil
}
