package swyftx

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

const (
	BaseURL = "https://api.swyftx.com.au/"
	// BaseURL = "https://private-anon-16c4713dbe-swyftx.apiary-mock.com/"
)

// Client holds the connection to swyftx and the api key and token for authentication
type Client struct {
	httpConn  *http.Client
	apiKey    string
	token     string
	userAgent string
	ctx       context.Context
}

type service struct {
	client *Client
}

// NewClientWithContext will create a new client with a specified context that can be used to
// interact with swyftx, if token is "" then a new token will be generated
func NewClientWithContext(ctx context.Context, apiKey, token string) (*Client, error) {
	client := &Client{
		token:  token,
		apiKey: apiKey,
		ctx:    ctx}

	// create http client for API
	cf := &tls.Config{Rand: rand.Reader}
	client.httpConn = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cf,
		},
	}

	client.userAgent = fmt.Sprintf("Manette/Alpha2 %s; Service", runtime.GOOS)

	if token == "" {
		var err error
		client.token, err = client.Authentication().Refresh()
		if err != nil {
			return nil, fmt.Errorf("could not generate a token: %s", err.Error())
		}
	}

	return client, nil
}

// NewClient will create a new client that can be used to interact with swyftx
// If token is "" then a new token will be generated
func NewClient(apiKey, token string) (*Client, error) {
	return NewClientWithContext(context.Background(), apiKey, token)
}

// NewRequest will create a new request that can be sent to the swyftx
func (c *Client) NewRequest(method, url string, body interface{}) (req *http.Request, err error) {
	var buf bytes.Buffer
	if body != nil {
		if err = json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("could not encode body of request: %s", err.Error())
		}
	}

	req, err = http.NewRequestWithContext(c.ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Add("Authorization", buildString("Bearer ", c.token))
	}
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// Do will do a request for the swyftx API and unmarshal the response into v
func (c *Client) Do(req *http.Request, v interface{}) (resp *http.Response, err error) {
	resp, err = c.httpConn.Do(req)
	if err != nil {
		return nil, err
	}

	var body *bytes.Buffer
	body, err = copyReadCloser(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not copy response body: %s", err.Error())
	}

	/*f, err := os.Create("debug.json")
	if err != nil {
		return nil, err
	}

	f.Write(body.Bytes())*/
	if resp.StatusCode >= http.StatusBadRequest {
		var errResp struct {
			Error Error `json:"error"`
		}
		if err = decodeJSON(body, &errResp); err != nil {
			return resp, fmt.Errorf("could not decode error: %s", err)
		}
		return resp, &errResp.Error
	}

	if err = decodeJSON(body, v); err != nil {
		return resp, fmt.Errorf("could not decode response: %s", err.Error())
	}

	return resp, nil
}

// Request will send a request to swyftx and check the response for errors
func (c *Client) Request(method, path string, body, v interface{}) error {
	req, err := c.NewRequest(method, buildString(BaseURL, path), body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err.Error())
	}

	var resp *http.Response
	resp, err = c.Do(req, v)
	if err != nil {
		return fmt.Errorf("could not do request: %s", err.Error())
	}
	defer resp.Body.Close()

	return nil
}

// Get http request to the Swyftx api
func (c *Client) Get(path string, v interface{}) error {
	return c.Request(http.MethodGet, path, nil, v)
}

// Post http request to the Swyftx api
func (c *Client) Post(path string, body, v interface{}) error {
	return c.Request(http.MethodPost, path, body, v)
}

// Delete http request to the Swyftx api
func (c *Client) Delete(path string) error {
	return c.Request(http.MethodDelete, path, nil, nil)
}

// Version of the Swyftx api
func (c *Client) Version() (string, error) {
	var version struct {
		Version string `json:"version"`
	}
	if err := c.Get("info/", &version); err != nil {
		return "", nil
	}

	return version.Version, nil
}

// WithContext will update the clients context to the one provided
func (c *Client) WithContext(ctx context.Context) *Client {
	c2 := new(Client)
	if ctx == nil {
		ctx = context.Background()
	}
	*c2 = *c
	c2.ctx = ctx
	return c2
}
