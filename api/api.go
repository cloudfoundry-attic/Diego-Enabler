package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PaginatedResponse struct {
	TotalPages int `json:"total_pages"`
}

type PageParser struct{}

func (p PageParser) Parse(body []byte) (PaginatedResponse, error) {
	var pages PaginatedResponse
	emptyPages := PaginatedResponse{}

	err := json.Unmarshal(body, &pages)
	if err != nil {
		return emptyPages, err
	}

	return pages, nil
}

type ApiClient struct {
	BaseUrl   *url.URL
	AuthToken string
}

func NewApiClient(rawurl string, authToken string) (*ApiClient, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	client := &ApiClient{
		BaseUrl:   u,
		AuthToken: authToken,
	}

	return client, nil
}

func (c *ApiClient) NewGetAppsRequest() (*http.Request, error) {
	req := &http.Request{
		Method: "GET",
		URL:    c.BaseUrl,
	}
	req.URL.Path = "/v2/apps"

	return req, nil
}

func (c *ApiClient) NewGetSpacesRequest() (*http.Request, error) {
	req := &http.Request{
		Method: "GET",
		URL:    c.BaseUrl,
	}
	req.URL.Path = "/v2/spaces"

	return req, nil
}

func (c *ApiClient) HandleFiltersAndParameters(next func() (*http.Request, error)) func(filter Filter, params map[string]interface{}) (*http.Request, error) {
	return func(filter Filter, params map[string]interface{}) (*http.Request, error) {
		req, err := next()
		if err != nil {
			return new(http.Request), err
		}

		req.URL.RawQuery = generateParams(filter, params).Encode()
		return req, nil
	}
}

func (c *ApiClient) Authorize(next func() (*http.Request, error)) func() (*http.Request, error) {
	return func() (*http.Request, error) {
		req, err := next()
		if err != nil {
			return new(http.Request), err
		}

		header := http.Header{}
		header.Set("Authorization", c.AuthToken)

		req.Header = header
		return req, nil
	}
}

func generateParams(filter Filter, params map[string]interface{}) url.Values {
	values := url.Values{}
	q := filter.ToFilterQueryParam()
	if q != "" {
		values.Set("q", q)
	}

	for k, v := range params {
		values.Set(k, fmt.Sprint(v))
	}
	return values
}
