package handlers

import (
	"net/http"

	"graphlabsts.core/internal/excel"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) GroupExcelTemplate(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.Repo.GetAllSubjectStrings()
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	bytes, err := excel.GetGroupTemplateFile(subjects)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=group_template.xlsx")
	w.Header().Set("Content-Type", "application/vnd.ms-excel")

	_, err = w.Write(bytes.Bytes())
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
