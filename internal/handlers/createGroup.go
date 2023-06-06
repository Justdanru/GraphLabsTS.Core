package handlers

import (
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

	err = excel.ReadGroupFile(file)
}
