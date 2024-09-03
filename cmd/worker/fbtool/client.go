package fbtool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	httpClient *http.Client
	APIKey     string
}

func NewClient(token string) *Client {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	dc := &Client{
		httpClient: client,
		APIKey:     token,
	}

	return dc
}

type StatisticsResponse struct {
	Data         []StatisticsAccount `json:"data"`
	RequestsLeft int                 `json:"requestsLeft"`

	// optional
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type StatisticsAccount struct {
	AccountID string `json:"account_id"`
	Currency  string `json:"currency"`

	Ads struct {
		Data []StatisticsAccountAd `json:"data"`
	} `json:"ads"`
}

type StatisticsAccountAd struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	EffectiveStatus string `json:"effective_status"`
	Status          string `json:"status"`

	Insights struct {
		Data []StatisticsAccountAdInsight `json:"data"`
	} `json:"insights"`
}

type StatisticsAccountAdInsight struct {
	Impressions int       `json:"impressions"`
	Clicks      int       `json:"clicks"`
	Spend       float64   `json:"spend"`
	Date        time.Time `json:"date_start"`
}

func (a *StatisticsAccountAdInsight) UnmarshalJSON(data []byte) error {
	type auxStatisticsAccountAdInsight struct {
		Impressions string `json:"impressions"`
		Clicks      string `json:"clicks"`
		Spend       string `json:"spend"`
		Date        string `json:"date_start"`
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

	if aux.Date != "" {
		date, err := time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return err
		}

		a.Date = date
	}

	return nil
}

func (c *Client) GetStatistics(account int, start, end time.Time) (*StatisticsResponse, error) {
	end_date := end.Format("2006-01-02")
	start_date := start.Format("2006-01-02")

	url := "https://fbtool.pro/api/get-statistics?key=%s&account=%d&dates=%s+-+%s&byDay=1"
	reqURL := fmt.Sprintf(url, c.APIKey, account, start_date, end_date)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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
	Data         map[string]*Account `json:"-"`
	RequestsLeft int                 `json:"requestsLeft"`
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

func (c *Client) GetAccounts() (*AccountsResponse, error) {
	url := "https://fbtool.pro/api/get-accounts?key=%s"
	reqURL := fmt.Sprintf(url, c.APIKey)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &AccountsResponse{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	accounts := make(map[string]*Account)
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}
	res.Data = accounts

	return res, nil
}
