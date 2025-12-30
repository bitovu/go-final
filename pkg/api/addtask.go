package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"final/pkg/db"
)

type AddTaskResp struct {
	ID string `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	now := time.Now()
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "request body parse error",
		})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "title is required",
		})
		return
	}

	if err := checkDate(&task, now); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "date is required",
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "task add error",
		})
		return
	}

	writeJSON(w, http.StatusOK, AddTaskResp{
		ID: strconv.FormatInt(id, 10),
	})
}
