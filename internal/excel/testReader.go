package excel

import (
	"errors"
	"graphlabsts/internal/types"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// TODO Обернуть ошибку
func getQuestionsInFile(fileWithTest *excelize.File) ([]types.Question, error) {
	result := []types.Question{}
	i := 3
	for {
		value, err := fileWithTest.GetCellValue("Вопросы", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		if value == "" {
			break
		}
		question := types.Question{}
		switch value {
		case "ОдинОтвет":
			question.QuestionTypeCode = types.OneAnswerQuestionCode
		case "МножОтвет":
			question.QuestionTypeCode = types.MultAnswerQuestionCode
		case "ЛожьИстина":
			question.QuestionTypeCode = types.TrueFalseQuestionCode
		case "Упорядоч":
			question.QuestionTypeCode = types.OrderQuestionCode
		case "Сопост":
			question.QuestionTypeCode = types.MatchQuestionCode
		default:
			return nil, errors.New("unknown type of question in B" + strconv.Itoa(i) + ": " + value)
		}
		// TODO Валидировать данные
		value, err = fileWithTest.GetCellValue("Вопросы", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		question.Text = value
		// TODO Валидировать данные
		value, err = fileWithTest.GetCellValue("Вопросы", "D"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		question.AnswerOptions = value
		value, err = fileWithTest.GetCellValue("Вопросы", "E"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		// TODO Валидировать данные
		question.CorrectAnswer = value
		value, err = fileWithTest.GetCellValue("Вопросы", "F"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		s := strings.Split(value, ";")
		for _, v := range s {
			// TODO Валидировать данные
			question.Terms = append(question.Terms, strings.TrimSpace(v))
		}
		// TODO Валидировать данные
		value, err = fileWithTest.GetCellValue("Вопросы", "G"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		question.Score, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, question)
		i++
	}
	return result, nil
}

// TODO Сделать конфигурационный файл с названиями листов и указанием ячеек
// TODO Возможно, сделать обёртку для возвращаемой ошибки
func readTestInFile(fileWithTest *excelize.File, test *types.Test) error {
	// TODO Валидировать данные
	value, err := fileWithTest.GetCellValue("Основная информация", "C2")
	if err != nil {
		return err
	}
	test.Title = value
	// TODO Валидировать данные
	value, err = fileWithTest.GetCellValue("Основная информация", "C3")
	if err != nil {
		return err
	}
	test.Subject = value
	value, err = fileWithTest.GetCellValue("Основная информация", "C4")
	if err != nil {
		return err
	}
	s := strings.Split(value, ";")
	for _, v := range s {
		// TODO Валидировать данные
		test.Groups = append(test.Groups, strings.TrimSpace(v))
	}
	value, err = fileWithTest.GetCellValue("Основная информация", "C5")
	if err != nil {
		return err
	}
	// TODO Валидировать данные
	test.BeginTime = value
	value, err = fileWithTest.GetCellValue("Основная информация", "C6")
	if err != nil {
		return err
	}
	// TODO Валидировать данные
	test.EndTime = value
	// TODO Сделать считывание количества вопросов
	value, err = fileWithTest.GetCellValue("Основная информация", "C8")
	if err != nil {
		return err
	}
	if value == "Адаптивный" {
		test.IsAdaptive = true
	} else if value == "Простой" {
		test.IsAdaptive = false
	} else {
		return errors.New("unknown test mode")
	}
	i := 3
	for {
		value, err := fileWithTest.GetCellValue("Термины и понятия", "B"+strconv.Itoa(i))
		if err != nil {
			return err
		}
		if value == "" {
			break
		}
		// TODO Валидировать данные
		for k := range test.Terms {
			if test.Terms[k] == value {
				return errors.New("duplicated term '" + value + "'")
			}
		}
		test.Terms = append(test.Terms, value)
		i++
	}
	test.Questions, err = getQuestionsInFile(fileWithTest)
	if err != nil {
		return err
	}
	return nil
}

// TODO Сменить аргументы на один цельный relativePathToFile
func GetTestFromFile(relativePathToFile string) (test *types.Test, err error) {
	wrapErr := errors.New("GetTestFromFile(" + relativePathToFile + ")")
	fileWithTest, inErr := excelize.OpenFile(relativePathToFile)
	if inErr != nil {
		return nil, errors.Join(wrapErr, inErr)
	}
	defer func() {
		inErr := fileWithTest.Close()
		if inErr != nil {
			err = errors.Join(wrapErr, fileWithTest.Close())
		}
	}()
	test = &types.Test{}
	err = readTestInFile(fileWithTest, test)
	return test, err
}
