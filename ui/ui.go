package ui

import (
	"fmt"
	"strings"
)

type Runtime string

const (
	DEA   Runtime = "DEA"
	Diego Runtime = "Diego"
)

func (r Runtime) String() string {
	switch strings.ToLower(string(r)) {
	case strings.ToLower(string(DEA)):
		return string(DEA)
	case strings.ToLower(string(Diego)):
		return string(Diego)
	default:
		return string(r)
	}
}

func (r Runtime) Flip() Runtime {
	if r == DEA {
		return Diego
	}
	return DEA
}

func ParseRuntime(runtime string) (Runtime, error) {
	switch strings.ToLower(runtime) {
	case "dea":
		return DEA, nil
	case "diego":
		return Diego, nil
	default:
		return "", fmt.Errorf("unknown runtime %s", runtime)
	}
}

type ApplicationPrinter interface {
	Name() string
	Organization() string
	Space() string
}
