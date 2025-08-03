package netbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ayofash1/netbox-checker/internal/rules"
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
			Timeout: 15 * time.Second,
		},
	}
}

// FetchDevices fetches all devices and their related data
func (c *Client) FetchDevices() ([]rules.Device, error) {
	var devices []rules.Device
	nextURL := fmt.Sprintf("%s/api/dcim/devices/?limit=100", c.BaseURL)

	for nextURL != "" {
		req, _ := http.NewRequest("GET", nextURL, nil)
		req.Header.Set("Authorization", "Token "+c.Token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch devices: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("NetBox error %d: %s", resp.StatusCode, string(body))
		}

		var page DevicePage
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("json error: %w", err)
		}

		for _, d := range page.Results {
			interfaces, err := c.fetchInterfaces(d.ID)
			if err != nil {
				return nil, fmt.Errorf("fetching interfaces for device %s: %w", d.Name, err)
			}

			ruleDev := MapToRuleDevice(d)
			ruleDev.Interfaces = interfaces
			devices = append(devices, ruleDev)
		}

		// Prepare next URL
		if page.Next != nil {
			parsed, _ := url.Parse(*page.Next)
			nextURL = c.BaseURL + parsed.RequestURI()
		} else {
			nextURL = ""
		}
	}

	return devices, nil
}

func (c *Client) fetchInterfaces(deviceID int) ([]string, error) {
	url := fmt.Sprintf("%s/api/dcim/interfaces/?device_id=%d&limit=100", c.BaseURL, deviceID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Token "+c.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interfaces: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("NetBox returned %d: %s", resp.StatusCode, string(body))
	}

	var page InterfacePage
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("error decoding interface JSON: %w", err)
	}

	var names []string
	for _, iface := range page.Results {
		names = append(names, iface.Name)
	}
	return names, nil
}
