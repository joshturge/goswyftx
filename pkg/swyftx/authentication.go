package swyftx

// Scope contains permission information for a particular API key action
type Scope struct {
	Display     string `json:"display"`
	Description string `json:"desc"`
	Key         string `json:"key"`
	State       int    `json:"state"`
}

// AppScope contains information about an API keys scope
type AppScope struct {
	ReadAccount   Scope `json:"app.account.read"`
	WithdrawFunds Scope `json:"app.funds.withdraw"`
	DeleteOrders  Scope `json:"app.orders.delete"`
}

// Key contains information about an API key
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

// Refresh will regenerate a new access token (JWT token) given an api key
func (as *AuthService) Refresh(apiKey string) (string, error) {
	var (
		token struct {
			Token string `json:"accessToken"`
		}
		body struct {
			APIKey string `json:"apiKey"`
		}
	)
	body.APIKey = apiKey

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

// Scope will get the scope of permmissions for the current API key
func (as *AuthService) Scope() (*AppScope, error) {
	var appScope AppScope
	if err := as.client.Get("user/apiKeys/scope/", &appScope); err != nil {
		return nil, err
	}

	return &appScope, nil
}

// Keys will get all the keys available for the current swyftx account
func (as *AuthService) Keys() ([]*Key, error) {
	var keys []*Key
	if err := as.client.Get("user/apiKeys/", &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

// RevokeKey will revoke an API key for the current swyftx account
// Returns the status of that action
func (as *AuthService) RevokeKey(apiKey string) (string, error) {
	var status struct {
		Status string `json:"status"`
	}
	if err := as.client.Post("user/apiKeys/revoke/", &apiKey, &status); err != nil {
		return "", err
	}

	return status.Status, nil
}

// RevokeAllKeys will revoke all API keys for the current swyftx account
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
