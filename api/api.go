package api

import (
	"fmt"
	"net/http"
)

type PaginatedResponse struct {
	TotalPages int `json:"total_pages"`
}

type ApiClient struct{}

func (c ApiClient) NewGetAppsRequest(filters Filters, params map[string]interface{}) (*http.Request, error) {
	return &http.Request{}, nil
}

type Filters []Filter

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
