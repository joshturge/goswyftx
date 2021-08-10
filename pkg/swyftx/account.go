package swyftx

import "net/http"

// AccountService wraps account specific swyftx API endpoints
type AccountService service

// AccountSettings contains swyftx account specific settings
type AccountSettings struct {
	// Used only for Account Profile request, might be a bug with the api?
	FavouriteAssets struct {
		AssetID bool `json:"assetId,omitempty"`
	} `json:"favouriteAssets"`
	// Used only for updating account settings, might be a bug with the api?
	FavouriteAsset struct {
		AssetID   bool `json:"assetId,omitempty"`
		FavStatus bool `json:"favStatus,omitempty"`
	} `json:"favouriteAsset,omitempty"`
	AnalyticsOptOut    bool `json:"analyticsOptOut,omitempty"`
	ActiveAffiliation  bool `json:"activeAffil,omitempty"`
	DisableSMSRecovery bool `json:"disableSMSRecovery,omitempty"`
}

// AccountProfile containts swyftx account information
type AccountProfile struct {
	DOB  SwyftxTime `json:"dob,omitempty"`
	Name struct {
		First string `json:"first,omitempty"`
		Last  string `json:"last,omitempty"`
	} `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Currency struct {
		ID   int    `json:"id,omitempty"`
		Code string `json:"code,omitempty"`
	} `json:"currency,omitempty"`
	UserHash string `json:"user_hash,omitempty"`
	MetaData struct {
		MFAEnabled  bool `json:"mfa_enabled,omitempty"`
		MFAEnrolled bool `json:"mfa_enrolled,omitempty"`
	} `json:"metadata,omitempty"`
	// NOTE: what's this???
	UserSettings struct {
	} `json:"userSettings,omitempty"`
}

// AccountVerification contains swyftx account verification information
type AccountVerification struct {
	Status   string `json:"status,omitempty"`
	Email    string `json:"email,omitempty"`
	MFA      string `json:"mfa,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Identity string `json:"identity,omitempty"`
}

// AccountUserVerification contains a swyftx post verification status for either
// phone or email types.
type AccountUserVerification struct {
	Success       bool   `json:"success,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	PhoneVerified bool   `json:"phoneVerified,omitempty"`
	Msg           string `json:"msg,omitempty"`
}

// AccountAffiliation contains swyftx referral information
type AccountAffiliation struct {
	ReferralLink  string `json:"referral_link,omitempty"`
	ReferredUsers int    `json:"referred_users"`
}

// AccountBalance contains the balance of an asset
type AccountBalance struct {
	AssetID int `json:"assetId,omitempty"`
	AvailableBalance string `json:"availableBalance,omitempty"`
}

// AccountStatistics contains statistics for a swyftx account
type AccountStatistics struct {
	Orders    int     `json:"orders,omitempty"`
	Traded    float32 `json:"traded,omitempty"`
	Deposited float32 `json:"deposited,omitempty"`
	Withdrawn float32 `json:"withdrawn,omitempty"`
}

// AccountMilestones contains milestones for a swyftx account
type AccountMilestones struct {
	SignUp    bool `json:"signUp,omitempty"`
	Verified  bool `json:"verified,omitempty"`
	Deposit   bool `json:"deposit,omitempty"`
	Trade     bool `json:"trade,omitempty"`
	Refer     bool `json:"refer,omitempty"`
	Completed bool `json:"completed,omitempty"`
}

// Account returns an account service that can interact with swyftx API account
// endpoints
func (c *Client) Account() *AccountService {
	return (*AccountService)(&service{c})
}

// Profile will get a users profile
func (as *AccountService) Profile() (*AccountProfile, error) {
	var account struct {
		Profile AccountProfile `json:"profile"`
	}
	if err := as.client.Get("user/", &account); err != nil {
		return nil, err
	}

	return &account.Profile, nil
}

// UpdateSettings will update a users account settings
func (as *AccountService) UpdateSettings(accSett *AccountSettings) (*AccountProfile, error) {
	var (
		account struct {
			Profile AccountProfile `json:"profile"`
		}
		body struct {
			Data AccountSettings `json:"data"`
		}
	)
	body.Data = *accSett

	if err := as.client.Post("user/settings/", &body, &account); err != nil {
		return nil, err
	}

	return &account.Profile, nil
}

// Verification will get all user verification information
func (as *AccountService) VerificationInfo() (*AccountVerification, error) {
	var account struct {
		Verification AccountVerification `json:"verification"`
	}
	if err := as.client.Get("user/verification/", &account); err != nil {
		return nil, err
	}

	return &account.Verification, nil
}

// SaveGreenID will save GreenID verification info
func (as *AccountService) SaveGreenID(greenID string) error {
	var body struct {
		Verification struct {
			ID string `json:"id"`
		} `json:"verification"`
	}
	body.Verification.ID = greenID

	if err := as.client.Request(http.MethodGet, "user/verification/storeGreenId/", &body, nil); err != nil {
		return err
	}

	return nil
}

// StartEmailVerification will send an email to verify access to an email account
func (as *AccountService) StartEmailVerification() (*AccountUserVerification, error) {
	return as.startVerification("email", "")
}

// CheckEmailVerification will check the verification status of an email account
func (as *AccountService) CheckEmailVerification() (*AccountUserVerification, error) {
	return as.checkVerification("email", "")
}

// CheckPhoneVerification will send an SMS to a phone number containing a verification token
func (as *AccountService) CheckPhoneVerification(phone string) (*AccountUserVerification, error) {
	return as.checkVerification("phone", phone)
}

// StartPhoneVerification will try and verify access to a phone given a token
func (as *AccountService) StartPhoneVerification(token string) (*AccountUserVerification, error) {
	return as.startVerification("phone", token)
}

func (as *AccountService) startVerification(verifyType, token string) (*AccountUserVerification, error) {
	var accUserVerif AccountUserVerification
	if err := as.client.Post(buildString("user/verification/", verifyType, "/", token), nil,
		&accUserVerif); err != nil {
		return nil, err
	}

	return &accUserVerif, nil
}

func (as *AccountService) checkVerification(verifyType, phone string) (*AccountUserVerification, error) {
	var accUserVerif AccountUserVerification
	if err := as.client.Get(buildString("user/verification/", verifyType, "/", phone),
		&accUserVerif); err != nil {
		return nil, err
	}

	return &accUserVerif, nil
}

// Affiliation will get a user's affiliation link and statistics
func (as *AccountService) Affiliation() (*AccountAffiliation, error) {
	var accAffil AccountAffiliation
	if err := as.client.Get("user/affiliations/", &accAffil); err != nil {
		return nil, err
	}

	return &accAffil, nil
}

// Balances will get a user's account balances
func (as *AccountService) Balances() ([]*AccountBalance, error) {
	var balances []*AccountBalance
	if err := as.client.Get("user/balance/", &balances); err != nil {
		return nil, err
	}

	return balances, nil
}

// SetCurrency will update a user's default currency given the asset ID of the new currency
func (as *AccountService) SetCurrency(assetID int) (*AccountProfile, error) {
	var (
		body struct {
			Profile struct {
				DefaultAsset int `json:"defaultAsset"`
			} `json:"profile"`
		}
		account struct {
			Profile AccountProfile `json:"profile"`
		}
	)
	body.Profile.DefaultAsset = assetID

	if err := as.client.Post("user/currency/", &body, &account); err != nil {
		return nil, err
	}

	return &account.Profile, nil
}

// Statistics for the user's account and usage
func (as *AccountService) Statistics() (*AccountStatistics, error) {
	var stats AccountStatistics
	if err := as.client.Get("user/statistics/", &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// Progress shows the completion status of particular milestones for the user's account
func (as *AccountService) Progress() (*AccountMilestones, error) {
	var milstones AccountMilestones
	if err := as.client.Get("user/progress/", &milstones); err != nil {
		return nil, err
	}

	return &milstones, nil
}
