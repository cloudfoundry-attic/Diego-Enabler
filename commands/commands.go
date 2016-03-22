package commands

import "github.com/cloudfoundry-incubator/diego-enabler/api"

//go:generate counterfeiter . PaginatedRequester
type PaginatedRequester interface {
	Do(filter api.Filter, params map[string]interface{}) ([][]byte, error)
}
