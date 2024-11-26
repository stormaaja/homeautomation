package tools

import (
	"fmt"
	"io"
)

func ReadBodyString(readCloser *io.ReadCloser) (string, error) {
	var body, error = io.ReadAll(*readCloser)
	if error != nil {
		return "", error
	}
	return string(body), nil
}

func ReadBodyFloat(readCloser *io.ReadCloser) (float64, error) {
	bodyString, error := ReadBodyString(readCloser)

	if error != nil {
		return 0, error
	}
	var temperature float64

	error = nil

	_, error = fmt.Sscanf(bodyString, "%f", &temperature)

	return temperature, error
}
