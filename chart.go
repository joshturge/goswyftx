package goswyftx

import (
	"strconv"
	"time"
)

type ChartService service

type GetBarChartRequest struct {
	BaseAsset        string
	SecondaryAsset   string
	Resolution       string
	From             time.Time
	To               time.Time
	FirstDataRequest bool
}

type OCHLVT struct {
	Time   SwyftxTime `json:"time,omitempty"`
	Open   string     `json:"open,omitempty"`
	High   string     `json:"high,omitempty"`
	Low    string     `json:"low,omitempty"`
	Close  string     `json:"close,omitempty"`
	Volume int        `json:"volume,omitempty"`
}

type ChartAsset struct {
	BaseAsset  string `json:"baseAsset,omitempty"`
	Asset      string `json:"asset,omitempty"`
	Resolution string `json:"resolution,omitempty"`
}

type ChartSettings struct {
	SupportsSearch       bool `json:"supports_search,omitempty"`
	SupportsGroupRequest bool `json:"supports_group_request,omitempty"`
	SupportsMarks        bool `json:"supports_marks,omitempty"`
	// Warning unknown typ of array, using string as default
	Exchanges []string `json:"exchanges,omitempty"`
	// Warning unknown typ of array, using string as default
	SymbolsTypes         []string `json:"symbols_types,omitempty"`
	SupportedResolutions string   `json:"supported_resolutions"`
}

type ChartResolveSymbol struct {
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	Type                 string `json:"type,omitempty"`
	Session              string `json:"session,omitempty"`
	Exchange             string `json:"exchange,omitempty"`
	ListedExchange       string `json:"listed_exchange,omitempty"`
	Timezone             string `json:"timezone,omitempty"`
	MinMov               int    `json:"minmov,omitempty"`
	PriceScale           int    `json:"pricescale,omitempty"`
	MinMive2             int    `json:"minmive2,omitempty"`
	HasIntraday          bool   `json:"has_intraday,omitempty"`
	SupportedResolutions string `json:"supported_resolutions,omitempty"`
	DataStatus           string `json:"data_status,omitempty"`
}

// Chart will return a chart service that can interact with the swyftx api
func (c *Client) Chart() *ChartService {
	return (*ChartService)(&service{c})
}

// Bar chart that contains pricing ticks for asset pair
func (cs *ChartService) Bar(cRequest *GetBarChartRequest) ([]*OCHLVT, error) {
	uri := buildString("charts/getBars/",
		cRequest.BaseAsset, "/",
		cRequest.SecondaryAsset, "/",
		cRequest.Resolution, "/",
		"?from=", strconv.FormatInt((cRequest.From.UnixNano()/int64(time.Millisecond)), 10),
		"&to=", strconv.FormatInt((cRequest.To.UnixNano()/int64(time.Millisecond)), 10),
		"&firstDataRequest=", strconv.FormatBool(cRequest.FirstDataRequest),
	)

	var barCharts []*OCHLVT
	if err := cs.client.Get(uri, &barCharts); err != nil {
		return nil, err
	}

	return barCharts, nil
}

// LatestBar will return the latest bar for an asset/s
func (cs *ChartService) LatestBar(cAssets ...ChartAsset) ([]*OCHLVT, error) {
	var (
		body      []ChartAsset
		barCharts []*OCHLVT
	)
	body = cAssets

	if err := cs.client.Post("charts/getLatestBar/", &body, &barCharts); err != nil {
		return nil, err
	}

	return barCharts, nil
}

// Settings will get chart settings
func (cs *ChartService) Settings() (*ChartSettings, error) {
	var cSettings ChartSettings
	if err := cs.client.Get("charts/settings", &cSettings); err != nil {
		return nil, err
	}

	return &cSettings, nil
}

// ResolveSymbols will get return a resolve symbol for a crypto pair
func (cs *ChartService) ResolveSymbols(baseAsset, secondaryAsset int) (*ChartResolveSymbol, error) {
	var resolveSymbol ChartResolveSymbol
	if err := cs.client.Get(buildString("charts/resolveSymbol/",
		strconv.Itoa(baseAsset),
		strconv.Itoa(secondaryAsset)), &resolveSymbol); err != nil {
		return nil, err
	}

	return &resolveSymbol, nil
}
