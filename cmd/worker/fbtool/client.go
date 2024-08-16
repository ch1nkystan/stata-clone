package fbtool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const host = "https://fbtool.pro"

type Client struct {
	HTTPClient *http.Client
	APIKey     string
}

func NewClient(token string) *Client {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	dc := &Client{
		HTTPClient: client,
		APIKey:     token,
	}

	return dc
}

type StatisticsResponse struct {
	Data []StatisticsAccount `json:"data"`

	// optional
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type StatisticsAccount struct {
	AccountID string `json:"account_id"`
	Currency  string `json:"currency"`

	Ads StatisticsAccountAds `json:"ads"`
}

type StatisticsAccountAds struct {
	Data []StatisticsAccountAd `json:"data"`
}

type StatisticsAccountAd struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	EffectiveStatus string `json:"effective_status"`
	Status          string `json:"status"`

	Insights StatisticsAccountAdInsights `json:"insights"`
}

type StatisticsAccountAdInsights struct {
	Data []StatisticsAccountAdInsight `json:"data"`
}

type StatisticsAccountAdInsight struct {
	Impressions int     `json:"impressions"`
	Clicks      int     `json:"clicks"`
	Spend       float64 `json:"spend"`
}

func (a *StatisticsAccountAdInsight) UnmarshalJSON(data []byte) error {
	type auxStatisticsAccountAdInsight struct {
		Impressions string `json:"impressions"`
		Clicks      string `json:"clicks"`
		Spend       string `json:"spend"`
	}

	aux := &auxStatisticsAccountAdInsight{}
	if err := json.Unmarshal(data, aux); err != nil {
		return nil
	}

	if aux.Impressions != "" {
		impressions, err := strconv.Atoi(aux.Impressions)
		if err != nil {
			return err
		}

		a.Impressions = impressions
	}

	if aux.Clicks != "" {
		clicks, err := strconv.Atoi(aux.Clicks)
		if err != nil {
			return err
		}

		a.Clicks = clicks
	}

	if aux.Spend != "" {
		spend, err := strconv.ParseFloat(aux.Spend, 64)
		if err != nil {
			return err
		}

		a.Spend = spend
	}

	return nil
}

func (c *Client) GetStatistics(account int, day time.Time) (*StatisticsResponse, error) {
	date := day.Format("2006-01-02")

	u := fmt.Sprintf("%s/api/get-statistics?key=%s&account=%d&dates=%s+-+%s",
		host, c.APIKey, account, date, date)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get statistics, invalid status code: %s", resp.Status)
	}

	defer resp.Body.Close()
	res := &StatisticsResponse{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return res, nil
}

type AccountsResponse struct {
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (a *Account) UnmarshalJSON(data []byte) error {
	type auxAccount struct {
		ID   string
		Name string
	}

	aux := &auxAccount{}
	if err := json.Unmarshal(data, aux); err != nil {
		return nil
	}

	id, err := strconv.Atoi(aux.ID)
	if err != nil {
		return err
	}

	a.ID = id
	a.Name = aux.Name

	return nil
}

func (c *Client) GetAccounts() (map[string]*Account, error) {
	u := fmt.Sprintf("%s/api/get-accounts?key=%s", host, c.APIKey)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	res := make(map[string]*Account)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
