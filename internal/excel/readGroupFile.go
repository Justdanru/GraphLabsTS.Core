package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/validators"
)

func ReadGroupFile(source io.Reader) (*models.GroupFromFile, error) {
	s, err := excelize.OpenReader(source)
	if err != nil {
		return nil, ErrOpenFile
	}

	group := &models.GroupFromFile{}

	var isValid bool

	group.Name, err = s.GetCellValue("Название группы", "A2")
	if err != nil {
		return nil, ErrReadingFile
	}
	isValid, err = validators.IsGroupNameValid(group.Name)
	if err != nil {
		return nil, ErrReadingFile
	}
	if !isValid {
		return nil, ErrWrongGroupNameFormat
	}

	i := 2
	for {
		subject, err := s.GetCellValue("Дисциплины группы", "A"+fmt.Sprint(i))
		if err != nil {
			return nil, ErrReadingFile
		}

		if (subject == "") && (i == 2) {
			return nil, ErrNoGroupSubjects
		}
		if (subject == "") && (i > 2) {
			break
		}

		isValid, err = validators.IsSubjectTitleValid(subject)
		if err != nil {
			return nil, ErrReadingFile
		}
		if !isValid {
			return nil, ErrWrongSubjectTitleFormat
		}

		i++
		group.Subjects = append(group.Subjects, subject)
	}

	i = 2
	for {
		student := &models.StudentFromFile{}

		student.Surname, err = s.GetCellValue("Состав группы", "A"+fmt.Sprint(i))
		if err != nil {
			return nil, ErrReadingFile
		}

		if (student.Surname == "") && (i == 2) {
			return nil, ErrNoStudents
		}
		if (student.Surname == "") && (i > 2) {
			break
		}

		isValid, err = validators.IsUserSurnameValid(student.Surname)
		if err != nil {
			return nil, ErrReadingFile
		}
		if !isValid {
			return nil, ErrWrongSurnameFormat
		}

		student.Name, err = s.GetCellValue("Состав группы", "B"+fmt.Sprint(i))
		if err != nil {
			return nil, ErrReadingFile
		}

		isValid, err = validators.IsUserNameValid(student.Name)
		if err != nil {
			return nil, ErrReadingFile
		}
		if !isValid {
			return nil, ErrWrongNameFormat
		}

		student.LastName, err = s.GetCellValue("Состав группы", "C"+fmt.Sprint(i))
		if err != nil {
			return nil, ErrReadingFile
		}

		isValid, err = validators.IsUserLastNameValid(student.LastName)
		if err != nil {
			return nil, ErrReadingFile
		}
		if !isValid {
			return nil, ErrWrongLastNameFormat
		}

		student.TelegramUsername, err = s.GetCellValue("Состав группы", "D"+fmt.Sprint(i))
		if err != nil {
			return nil, ErrReadingFile
		}

		isValid, err = validators.IsTelegramUsernameValid(student.TelegramUsername)
		if err != nil {
			return nil, ErrReadingFile
		}
		if !isValid {
			return nil, ErrWrongTelegramUsernameFormat
		}

		i++
		group.Students = append(group.Students, student)

	}

	return group, nil
}
