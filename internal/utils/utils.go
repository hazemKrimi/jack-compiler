package utils

import (
	"errors"
	"io"
)

func CheckReaderError(err error) (bool, bool) {
	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, true
		}

		return true, false
	}

	return false, false
}
