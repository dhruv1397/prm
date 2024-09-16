package util

import "fmt"

func FormatErrors(errs []error) string {
	errStrings := ""
	for _, err := range errs {
		errStrings += fmt.Sprintf("%v\n", err)
	}
	return errStrings
}
