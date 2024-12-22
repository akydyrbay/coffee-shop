package dal

import (
	"errors"
	"os"
)

func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(path)
	}
	return file, err
}
