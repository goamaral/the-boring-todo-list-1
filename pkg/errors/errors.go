package errors

import (
	"fmt"
)

func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %s", err.Error(), msg)
}

func Wrapf(err error, format string, a ...any) error {
	return fmt.Errorf("%s: %s", err.Error(), fmt.Sprintf(format, a...))
}

func NewParseError(itemName string, itemValue any) error {
	return fmt.Errorf("failed to parse %s (%+v)", itemName, itemValue)
}
