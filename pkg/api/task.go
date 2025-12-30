package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
	"strings"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "unsupported method",
		})
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id is required",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	now := time.Now()

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "decoding error",
		})
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id is required",
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
			Error: err.Error(),
		})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id is required",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "task not found",
		})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: "delete error",
			})
			return
		}

		writeJSON(w, http.StatusOK, nil)
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "next date error",
		})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "update error",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "id is required",
		})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}
