package goswyftx

import "strconv"

type AddressDetails struct {
}

type Address struct {
	ID      int    `json:"id,omitempty"`
	Code    string `json:"code,omitempty"`
	Details struct {
		Address string `json:"address,omitempty"`
		DestTag string `json:"dest_tag,omitempty"`
	} `json:"address_details,omitempty"`
	Time SwyftxTime `json:"time,omitempty"`
	Name string     `json:"name,omitempty"`
	Type string     `json:"type,omitempty"`
}

type BSBStatus struct {
	// Duration in milliseconds
	Duration          int    `json:"durationMs,omitempty"`
	Status            string `json:"status,omitempty"`
	StatusDescription string `json:"statusDescription,omitempty"`
	Address           string `json:"address,omitempty"`
	BankCode          string `json:"bankCode,omitempty"`
	BSB               string `json:"bsb,omitempty"`
	City              string `json:"city,omitempty"`
	Closed            bool   `json:"closed,omitempty"`
	PostCode          string `json:"postCode,omitempty"`
	State             string `json:"state,omitempty"`
}

// AddressService holds methods that can interact with Swyftx address endpoints
type AddressService struct {
	service
	assetCode string
}

// Address will create a new Address service that can interact with the Swyftx addresses endpoints
// The asset code is required for the Deposit, Withdraw and CheckDeposit endpoints
func (c *Client) Address(assetCode ...string) *AddressService {
	if len(assetCode) == 0 {
		assetCode[0] = ""
	}

	return &AddressService{service{c}, assetCode[0]}
}

// Create will create a new address for a specific asset and return the newly created address
func (as *AddressService) Create(name string) (*Address, error) {
	if isEmptyStr(as.assetCode) {
		return nil, errAssetCode
	}

	var (
		addresses []*Address
		body      struct {
			Address struct {
				Name string `json:"name"`
			} `json:"address"`
		}
	)
	body.Address.Name = name

	if err := as.client.Post(buildString("address/deposit/", as.assetCode), &body, &addresses); err != nil {
		return nil, err
	}

	return addresses[0], nil
}

// GetActive will get all active addresses for an asset
func (as *AddressService) GetActive() ([]*Address, error) {
	return as.getAddresses("deposit")
}

// GetSaved will get all saved addresses for an asset
func (as *AddressService) GetSaved() ([]*Address, error) {
	return as.getAddresses("withdraw")
}

func (as *AddressService) getAddresses(fiatType string) ([]*Address, error) {
	if isEmptyStr(as.assetCode) {
		return nil, errAssetCode
	}

	var addresses []*Address
	if err := as.client.Get(buildString("address/", fiatType, "/", as.assetCode), &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

// Remove will remove a withdrawal adddress given the id of the address
func (as *AddressService) Remove(addressID int) error {
	if err := as.client.Delete(buildString("address/withdraw/", strconv.Itoa(addressID))); err != nil {
		return err
	}

	return nil
}

// VerifyWithdrawal will verify a withdrawal given the verification token
func (as *AddressService) VerifyWithdrawal(token string) error {
	if err := as.client.Get(buildString("address/withdraw/verify/", token), nil); err != nil {
		return err
	}

	return nil
}

// VerifyBSB will verify a BSB number and send back the current status of that BSB
func (as *AddressService) VerifyBSB(bsb string) (*BSBStatus, error) {
	var bsbStatus BSBStatus
	if err := as.client.Get(buildString("address/withdraw/bsb-verify/", bsb), &bsbStatus); err != nil {
		return nil, err
	}

	return &bsbStatus, nil
}

// CheckDeposit check a deposit for an address given the address id
func (as *AddressService) CheckDeposit(addressID int) error {
	if isEmptyStr(as.assetCode) {
		return errAssetCode
	}

	if err := as.client.Get(buildString("address/check/", as.assetCode, "/", strconv.Itoa(addressID)),
		nil); err != nil {
		return err
	}

	return nil
}
