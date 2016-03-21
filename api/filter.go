package api

import (
	"fmt"
	"strings"
)

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

type InclusionFilter struct {
	Name   string
	Values []interface{}
}

func (f InclusionFilter) ToFilterQueryParam() string {
	var vals []string

	for _, v := range f.Values {
		vals = append(vals, fmt.Sprint(v))
	}

	return fmt.Sprintf("%s IN %v", f.Name, strings.Join(vals, ","))
}
