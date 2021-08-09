package swyftx

type Scope struct {
	Display     string `json:"display"`
	Description string `json:"desc"`
	Key         string `json:"key"`
	State       int    `json:"state"`
}

type AppScope struct {
	ReadAccount   Scope `json:"app.account.read"`
	WithdrawFunds Scope `json:"app.funds.withdraw"`
	DeleteOrders  Scope `json:"app.orders.delete"`
}

type Key struct {
	ID      string     `json:"id"`
	Label   string     `json:"label"`
	Scope   string     `json:"scope"`
	Created SwyftxTime `json:"created"`
}

// AuthService hold methods for the Authentication endpoints
// https://docs.swyftx.com.au/#/reference/authentication
type AuthService service

// Authentication will create a new AuthService instance
func (c *Client) Authentication() *AuthService {
	return (*AuthService)(&service{c})
}

// Refresh will regenerate a new access token (JWT token)
func (as *AuthService) Refresh() (string, error) {
	var (
		token struct {
			Token string `json:"accessToken"`
		}
		body struct {
			APIKey string `json:"apiKey"`
		}
	)
	body.APIKey = as.client.apiKey

	if err := as.client.Post("auth/refresh/", &body, &token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Logout will invalidate the current access token (JWT token)
func (as *AuthService) Logout() (bool, error) {
	var success struct {
		Success bool `json:"success"`
	}
	if err := as.client.Post("auth/logout/", nil, &success); err != nil {
		return success.Success, err
	}

	return success.Success, nil
}

// GetScope will get the scope of permmissions for an api key
func (as *AuthService) GetScope() (*AppScope, error) {
	var appScope AppScope
	if err := as.client.Get("user/apiKeys/scope/", &appScope); err != nil {
		return nil, err
	}

	return &appScope, nil
}

// GetKeys will get all the keys available to a user
func (as *AuthService) GetKeys() ([]*Key, error) {
	var keys []*Key
	if err := as.client.Get("user/apiKeys/", &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

// RevokeKey will revoke an api key
// Returns the status of that action
func (as *AuthService) RevokeKey() (string, error) {
	var status struct {
		Status string `json:"status"`
	}
	if err := as.client.Post("user/apiKeys/revoke/", &as.client.apiKey, &status); err != nil {
		return "", err
	}

	return status.Status, nil
}

// RevokeAllKeys will revoke all api keys for a user account
// Returns the status of that action
func (as *AuthService) RevokeAllKeys() (string, error) {
	var status struct {
		Status string `json:"status"`
	}
	if err := as.client.Post("user/apiKeys/revokeAll/", nil, &status); err != nil {
		return "", err
	}

	return status.Status, nil
}
