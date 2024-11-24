package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"taskmanager/models"
	"time"
)

var validate = validator.New()

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.WithError(err).Warn("Failed to decode request body")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ValidateStruct(task); err != nil {
		logger.WithError(err).Warn("Validation failed for task")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uint)
	task.UserID = userID

	result := DB.Create(&task)
	if result.Error != nil {
		logger.WithError(result.Error).Error("Failed to create task in database")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": task.ID,
	}).Info("Task created successfully")

	json.NewEncoder(w).Encode(task)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)

	var tasks []models.Task
	result := DB.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
		}).WithError(result.Error).Error("Failed to fetch tasks")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"task_count": len(tasks),
	}).Info("Fetched all tasks successfully")

	json.NewEncoder(w).Encode(tasks)
}

func GetTasksByPeriod(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)

	period := r.URL.Query().Get("period")

	var start, end time.Time
	now := time.Now()

	switch period {
	case "day":
		start = now.Truncate(24 * time.Hour)
		end = start.Add(24 * time.Hour)
	case "week":
		start = now.AddDate(0, 0, -int(now.Weekday()-1)).Truncate(24 * time.Hour)
		end = start.AddDate(0, 0, 7)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 1, 0)
	default:
		http.Error(w, "Invalid period. Use 'day', 'week', or 'month'.", http.StatusBadRequest)
		return
	}

	var tasks []models.Task
	result := DB.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).Find(&tasks)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var task models.Task
	result := DB.First(&task, id)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"task_id": id,
		}).Warn("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(uint)
	if task.UserID != userID {
		logger.WithFields(logrus.Fields{
			"task_id": id,
			"user_id": userID,
		}).Warn("Unauthorized access to task")
		http.Error(w, "Unauthorized access to task", http.StatusUnauthorized)
		return
	}

	logger.WithFields(logrus.Fields{
		"task_id": id,
		"user_id": userID,
	}).Info("Task retrieved successfully")

	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var task models.Task
	result := DB.First(&task, id)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"task_id": id,
		}).Warn("Task not found for update")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(uint)
	if task.UserID != userID {
		logger.WithFields(logrus.Fields{
			"task_id": id,
			"user_id": userID,
		}).Warn("Unauthorized update attempt")
		http.Error(w, "Unauthorized access to update task", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.WithError(err).Warn("Failed to decode request body during update")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ValidateStruct(task); err != nil {
		logger.WithError(err).Warn("Validation failed during task update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DB.Save(&task)
	logger.WithFields(logrus.Fields{
		"task_id": id,
		"user_id": userID,
	}).Info("Task updated successfully")

	json.NewEncoder(w).Encode(task)
}

// DeleteTask - Delete a task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var task models.Task
	result := DB.First(&task, id)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"task_id": id,
		}).Warn("Task not found for deletion")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(uint)
	if task.UserID != userID {
		logger.WithFields(logrus.Fields{
			"task_id": id,
			"user_id": userID,
		}).Warn("Unauthorized deletion attempt")
		http.Error(w, "Unauthorized access to delete task", http.StatusUnauthorized)
		return
	}

	DB.Delete(&task)
	logger.WithFields(logrus.Fields{
		"task_id": id,
		"user_id": userID,
	}).Info("Task deleted successfully")

	fmt.Fprintf(w, "Task deleted successfully")
}
