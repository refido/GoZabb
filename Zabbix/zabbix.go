package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Endpoint   string
	Auth       string
	HTTPClient *http.Client
}

type rpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error,omitempty"`
	ID int `json:"id"`
}

// New returns a client with a sane timeout.
func New(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login authenticates and stores the auth token on the client.
func (c *Client) Login(user, password string) error {
	req := rpcRequest{
		Jsonrpc: "2.0",
		Method:  "user.login",
		Params: map[string]string{
			"user":     user,
			"password": password,
		},
		ID: 1,
	}

	body, _ := json.Marshal(req)
	resp, err := c.HTTPClient.Post(c.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("login error: %s (%s)", res.Error.Message, res.Error.Data)
	}

	var token string
	if err := json.Unmarshal(res.Result, &token); err != nil {
		return err
	}

	c.Auth = token
	return nil
}

func (c *Client) Call(method string, params interface{}) (json.RawMessage, error) {
	req := rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Auth:    c.Auth,
		ID:      1,
	}

	body, _ := json.Marshal(req)
	resp, err := c.HTTPClient.Post(c.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, fmt.Errorf("zabbix error: %s (%s)", res.Error.Message, res.Error.Data)
	}
	return res.Result, nil
}

// Version returns the Zabbix API version (no auth required).
func (c *Client) Version() (string, error) {
	req := rpcRequest{
		Jsonrpc: "2.0",
		Method:  "apiinfo.version",
		Params:  map[string]interface{}{},
		ID:      1,
	}

	body, _ := json.Marshal(req)
	resp, err := c.HTTPClient.Post(c.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Error != nil {
		return "", fmt.Errorf("zabbix error: %s (%s)", res.Error.Message, res.Error.Data)
	}

	var version string
	if err := json.Unmarshal(res.Result, &version); err != nil {
		return "", err
	}
	return version, nil
}
