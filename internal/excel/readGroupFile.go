package excel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

func ReadGroupFile(source io.Reader) error {
	_, err := excelize.OpenReader(source)
	if err != nil {
		return ErrOpenFile
	}

	return nil
}
