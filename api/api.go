package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func (c *ApiClient) NewGetAppsRequest(filter Filter, params map[string]interface{}) (*http.Request, error) {
	req := &http.Request{
		Method: "GET",
		URL:    c.BaseUrl,
	}
	req.URL.Path = "/v2/apps"
	req.URL.RawQuery = generateParams(filter, params).Encode()

	header := http.Header{}
	header.Set("Authorization", c.AuthToken)

	req.Header = header

	return req, nil
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

type Filters []Filter

//go:generate counterfeiter . Filter
type Filter interface {
	ToFilterQueryParam() string
}

type EqualFilter struct {
	Name  string
	Value interface{}
}

func (f EqualFilter) ToFilterQueryParam() string {
	return fmt.Sprintf("%s:%v", f.Name, f.Value)
}

func (f Filters) ToFilterQueryParam() string {
	var filters []string

	for _, x := range f {
		filters = append(filters, x.ToFilterQueryParam())
	}

	return strings.Join(filters, ";")
}
