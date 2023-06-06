package excel

import "errors"

var (
	ErrOpenFile        = errors.New("error trying open file")
	ErrWritingToBuffer = errors.New("error writing file to buffer")
)
