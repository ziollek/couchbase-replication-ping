package utils

import "fmt"

func HandleError(format string, err error) {
	if err != nil {
		GetLogger().Fatalf(format, err)
	}
}

func ErrorString(err error) (result string) {
	if err != nil {
		result = fmt.Sprintf("%s", err)
	}
	return result
}
