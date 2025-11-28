package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (c *Client) Login(username, password string) (*LoginResponse, error) {
	body, err := c.doRequest("POST", "/api/v1/auth/login", LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	var resp LoginResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type BackupCreateRequest struct {
	DatabaseName string `json:"database_name"`
	Format       string `json:"format"`
	Compression  bool   `json:"compression"`
	BackupType   string `json:"backup_type"`
}

type BackupCreateResponse struct {
	RequestID uint   `json:"request_id"`
	Message   string `json:"message"`
}

func (c *Client) CreateBackup(req BackupCreateRequest) (*BackupCreateResponse, error) {
	body, err := c.doRequest("POST", "/api/v1/backups", req)
	if err != nil {
		return nil, err
	}

	var resp BackupCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type BackupStatus struct {
	ID           uint       `json:"id"`
	DatabaseName string     `json:"database_name"`
	Status       string     `json:"status"`
	DownloadURL  string     `json:"download_url,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	FileSize     int64      `json:"file_size"`
	ErrorMessage string     `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (c *Client) GetBackupStatus(requestID uint) (*BackupStatus, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/api/v1/backups/%d", requestID), nil)
	if err != nil {
		return nil, err
	}

	var status BackupStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) ListBackups() ([]BackupStatus, error) {
	body, err := c.doRequest("GET", "/api/v1/backups", nil)
	if err != nil {
		return nil, err
	}

	var backups []BackupStatus
	if err := json.Unmarshal(body, &backups); err != nil {
		return nil, err
	}

	return backups, nil
}
