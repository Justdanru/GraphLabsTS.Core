package handlers

import (
	"fmt"
	"net/http"

	"graphlabsts.core/internal/excel"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("fileInput")
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	group, err := excel.ReadGroupFile(file)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(group)
}
