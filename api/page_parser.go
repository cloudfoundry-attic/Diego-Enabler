package api

import "encoding/json"

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
