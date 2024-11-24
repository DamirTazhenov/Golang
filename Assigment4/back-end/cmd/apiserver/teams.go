package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"taskmanager/models"
	"time"
)

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)

	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Привязываем создателя команды как менеджера
	team.CreatorID = userID
	if err := DB.Create(&team).Error; err != nil {
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Добавляем создателя команды как участника с ролью "manager"
	member := models.TeamMember{
		UserID:   userID,
		TeamID:   team.ID,
		Role:     "manager",
		JoinedAt: time.Now(),
	}

	if err := DB.Create(&member).Error; err != nil {
		http.Error(w, "Failed to add team creator as member", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func AddTeamMember(w http.ResponseWriter, r *http.Request) {
	var member models.TeamMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uint)
	teamID, _ := strconv.Atoi(mux.Vars(r)["team_id"])

	// Проверяем, является ли вызывающий пользователь менеджером команды
	hasRole, err := CheckTeamRole(uint(teamID), userID, "manager")
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: only managers can add members", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пользователь еще не в команде
	if err := DB.Where("team_id = ? AND user_id = ?", member.TeamID, member.UserID).First(&models.TeamMember{}).Error; err == nil {
		http.Error(w, "User already in the team", http.StatusBadRequest)
		return
	}

	member.TeamID = uint(teamID)
	member.JoinedAt = time.Now()
	if err := DB.Create(&member).Error; err != nil {
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(member)
}

func GetTeamTasks(w http.ResponseWriter, r *http.Request) {
	teamID, _ := strconv.Atoi(mux.Vars(r)["team_id"])
	userID := r.Context().Value("user_id").(uint)

	// Проверяем, является ли пользователь членом команды
	if _, err := CheckTeamRole(uint(teamID), userID, "manager", "employee", "client"); err != nil {
		http.Error(w, "Unauthorized to view tasks", http.StatusUnauthorized)
		return
	}

	var tasks []models.Task
	if err := DB.Where("team_id = ?", teamID).Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func CheckTeamRole(teamID uint, userID uint, requiredRoles ...string) (bool, error) {
	var member models.TeamMember
	if err := DB.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error; err != nil {
		return false, err // Пользователь не найден в команде
	}

	// Проверяем, соответствует ли роль пользователя одной из требуемых ролей
	for _, role := range requiredRoles {
		if member.Role == role {
			return true, nil
		}
	}

	return false, nil
}
