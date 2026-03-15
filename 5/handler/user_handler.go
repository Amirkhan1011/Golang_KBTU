package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"practice5/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(q.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	var filter repository.UserFilter

	if idStr := q.Get("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}

	if name := q.Get("name"); name != "" {
		filter.Name = &name
	}
	if email := q.Get("email"); email != "" {
		filter.Email = &email
	}
	if gender := q.Get("gender"); gender != "" {
		filter.Gender = &gender
	}

	if bdFromStr := q.Get("birth_date_from"); bdFromStr != "" {
		if t, err := time.Parse("2006-01-02", bdFromStr); err == nil {
			filter.BirthDateFrom = &t
		}
	}
	if bdToStr := q.Get("birth_date_to"); bdToStr != "" {
		if t, err := time.Parse("2006-01-02", bdToStr); err == nil {
			filter.BirthDateTo = &t
		}
	}

	orderBy := q.Get("order_by")
	orderDir := q.Get("order_dir")

	resp, err := h.repo.GetPaginatedUsers(page, pageSize, filter, orderBy, orderDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	user1Str := q.Get("user1")
	user2Str := q.Get("user2")

	user1, err1 := strconv.Atoi(user1Str)
	user2, err2 := strconv.Atoi(user2Str)
	if err1 != nil || err2 != nil || user1 <= 0 || user2 <= 0 {
		http.Error(w, "user1 and user2 must be positive integers", http.StatusBadRequest)
		return
	}
	if user1 == user2 {
		http.Error(w, "user1 and user2 must be different", http.StatusBadRequest)
		return
	}

	friends, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
