package api

import (
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . RequestFactory
//TODO: Fix counterfeiter to find Filter correctly #NoFilter
type RequestFactory func(Filter, map[string]interface{}) (*http.Request, error)

//go:generate counterfeiter . CloudControllerClient
type CloudControllerClient interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate counterfeiter . PaginatedParser
type PaginatedParser interface {
	Parse([]byte) (PaginatedResponse, error)
}

type PaginatedRequester struct {
	RequestFactory RequestFactory
	Client         CloudControllerClient
	PageParser     PaginatedParser
}

func NewPaginatedRequester(cliConnection Connection, requestFactory RequestFactory) (*PaginatedRequester, error) {
	pageParser := PageParser{}

	httpClient, err := NewHttpClient(cliConnection)
	if err != nil {
		return nil, err
	}

	return &PaginatedRequester{
		RequestFactory: requestFactory,
		Client:         httpClient,
		PageParser:     pageParser,
	}, nil
}

func (p *PaginatedRequester) Do(filter Filter, params map[string]interface{}) ([][]byte, error) {
	var noBodies [][]byte

	req, err := p.RequestFactory(filter, params)
	if err != nil {
		return noBodies, err
	}

	var responseBodies [][]byte

	res, err := p.Client.Do(req)
	if err != nil {
		return noBodies, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return noBodies, err
	}

	responseBodies = append(responseBodies, body)

	paginatedRes, err := p.PageParser.Parse(body)
	if err != nil {
		return noBodies, err
	}
	for page := 2; page <= paginatedRes.TotalPages; page++ {
		// construct a new request with the current page
		params["page"] = page
		req, err := p.RequestFactory(filter, params)
		if err != nil {
			return noBodies, err
		}

		// perform the request
		res, err := p.Client.Do(req)
		if err != nil {
			return noBodies, err
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return noBodies, err
		}

		responseBodies = append(responseBodies, body)
	}

	return responseBodies, nil
}
