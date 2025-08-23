package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiEndpoint = "http://172.18.62.101/zabbix/api_jsonrpc.php"

type Client struct {
	Auth string
}

// Generic JSON-RPC request struct
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

// Authenticate to Zabbix
func Login(user, password string) (*Client, error) {
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
	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, fmt.Errorf("login error: %s", res.Error.Message)
	}

	var token string
	if err := json.Unmarshal(res.Result, &token); err != nil {
		return nil, err
	}

	return &Client{Auth: token}, nil
}

// Call Zabbix API
func (c *Client) Call(method string, params interface{}) (json.RawMessage, error) {
	req := rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Auth:    c.Auth,
		ID:      1,
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, fmt.Errorf("zabbix error: %s", res.Error.Message)
	}

	return res.Result, nil
}
