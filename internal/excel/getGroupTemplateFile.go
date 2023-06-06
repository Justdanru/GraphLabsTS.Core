package excel

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func GetGroupTemplateFile(subjects []string) (*bytes.Buffer, error) {
	// TODO Вынести в переменную окружения
	file, err := excelize.OpenFile("./excel_templates/group_template.xlsx")
	if err != nil {
		return nil, ErrOpenFile
	}
	defer file.Close()

	i := 2
	fmt.Println(subjects) //debug-output
	for _, subject := range subjects {
		file.SetCellValue("Доступные дисциплины", "A"+fmt.Sprint(i), subject)
		i++
	}

	fileBytes, err := file.WriteToBuffer()
	if err != nil {
		return nil, ErrWritingToBuffer
	}

	return fileBytes, nil
}
