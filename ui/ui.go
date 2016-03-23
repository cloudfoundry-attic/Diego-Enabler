package ui

import "strings"

type Runtime string

func (r Runtime) String() string {
	switch strings.ToLower(string(r)) {
	case strings.ToLower(DEA):
		return DEA
	case strings.ToLower(Diego):
		return Diego
	default:
		return string(r)
	}
}

const (
	DEA   = "DEA"
	Diego = "Diego"
)
