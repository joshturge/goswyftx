package swyftx

type LimitService service

type WithdrawLimit struct {
	Used            int `json:"used,omitempty"`
	Remaining       int `json:"remaining,omitempty"`
	Limit           int `json:"limit,omitempty"`
	RollingCycleHrs int `json:"rollingCycleHrs,omitempty"`
}

func (c *Client) Limit() *LimitService {
	return (*LimitService)(&service{c})
}

func (ls *LimitService) Withdrawal() (*WithdrawLimit, error) {
	var withLimit WithdrawLimit
	if err := ls.client.Get("limits/withdrawal/", &withLimit); err != nil {
		return nil, err
	}

	return &withLimit, nil
}
