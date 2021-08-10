package swyftx

// LimitService contains a method which can interact with the swyftx API limit
// endpoint
type LimitService service

// WithdrawalLimit contains information about a withdrawal limit order
type WithdrawLimit struct {
	Used            int `json:"used,omitempty"`
	Remaining       int `json:"remaining,omitempty"`
	Limit           int `json:"limit,omitempty"`
	RollingCycleHrs int `json:"rollingCycleHrs,omitempty"`
}

// Limit will create a new limit service which can be used to interact with
// swyftx API limit endpoints
func (c *Client) Limit() *LimitService {
	return (*LimitService)(&service{c})
}

// Withdrawal will get the withdrawal limit for a swyftx user
func (ls *LimitService) Withdrawal() (*WithdrawLimit, error) {
	var withLimit WithdrawLimit
	if err := ls.client.Get("limits/withdrawal/", &withLimit); err != nil {
		return nil, err
	}

	return &withLimit, nil
}
