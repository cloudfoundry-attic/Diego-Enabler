package flaghelpers

import (
	"fmt"
	"strconv"
)

type ParallelFlag struct {
	Value int
}

func (flag *ParallelFlag) UnmarshalFlag(value string) error {
	val, _ := strconv.Atoi(value)
	if val <= 0 || val > 100 {
		return InvalidParallelValueError{PassedValue: value}
	}

	flag.Value = val
	return nil
}

type InvalidParallelValueError struct {
	PassedValue string
}

func (e InvalidParallelValueError) Error() string {
	return fmt.Sprintf(
		"Invalid maximum apps in flight: %s\nValue for MAX_IN_FLIGHT must be an integer between 1 and 100",
		e.PassedValue,
	)
}
