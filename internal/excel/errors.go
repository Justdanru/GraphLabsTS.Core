package excel

import "errors"

var (
	ErrOpenFile                    = errors.New("error trying open file")
	ErrReadingFile                 = errors.New("error while reading file")
	ErrWritingToBuffer             = errors.New("error writing file to buffer")
	ErrWrongGroupNameFormat        = errors.New("wrong group name format")
	ErrWrongSubjectTitleFormat     = errors.New("wrong subject title format")
	ErrWrongNameFormat             = errors.New("wrong user name format")
	ErrWrongSurnameFormat          = errors.New("wrong user surname format")
	ErrWrongLastNameFormat         = errors.New("wrong user last name format")
	ErrWrongTelegramUsernameFormat = errors.New("wrong telegram username format")
	ErrNoGroupSubjects             = errors.New("no subjects connected to group")
	ErrNoStudents                  = errors.New("no students in group")
)
