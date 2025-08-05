package handlers

import (
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
	"github.com/gorilla/mux"
)

// CreateUser handles POST /api/v1/users
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.UserCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	user, err := h.services.User.CreateUser(r.Context(), req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" ||
			err.Error() == "user with firebase_uid "+req.FirebaseUID+" already exists" {
			h.writeErrorResponse(w, r, http.StatusConflict, "Conflict", err.Error())
			return
		}
		h.handleInternalError(w, r, err, "Failed to create user")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, user, "User created successfully")
}

// GetUser handles GET /api/v1/users/{id}
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	user, err := h.services.User.GetUserByID(r.Context(), id)
	if err != nil {
		if err.Error() == "user with id "+strconv.FormatInt(id, 10)+" not found" {
			h.handleNotFoundError(w, r, "User")
			return
		}
		h.handleInternalError(w, r, err, "Failed to get user")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, user, "User retrieved successfully")
}

// GetUserByUID handles GET /api/v1/users/uid/{firebase_uid}
func (h *Handlers) GetUserByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firebaseUID, exists := vars["firebase_uid"]
	if !exists || firebaseUID == "" {
		h.writeErrorResponse(w, r, http.StatusBadRequest, "Bad Request", "firebase_uid parameter is required")
		return
	}

	user, err := h.services.User.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		if err.Error() == "user with firebase_uid "+firebaseUID+" not found" {
			h.handleNotFoundError(w, r, "User")
			return
		}
		h.handleInternalError(w, r, err, "Failed to get user by Firebase UID")
		return
	}

	// Return only the user ID for security reasons (as in original code)
	response := map[string]int64{"id": user.ID}
	h.writeSuccessResponse(w, r, http.StatusOK, response, "User retrieved successfully")
}

// UpdateUser handles PUT /api/v1/users/{id}
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.UserUpdateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	user, err := h.services.User.UpdateUser(r.Context(), id, req)
	if err != nil {
		if err.Error() == "user not found" {
			h.handleNotFoundError(w, r, "User")
			return
		}
		h.handleInternalError(w, r, err, "Failed to update user")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, user, "User updated successfully")
}

// DeleteUser handles DELETE /api/v1/users/{id}
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	err = h.services.User.DeleteUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			h.handleNotFoundError(w, r, "User")
			return
		}
		h.handleInternalError(w, r, err, "Failed to delete user")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusNoContent, nil, "User deleted successfully")
}

// GetUserProjects handles GET /api/v1/users/{id}/projects
func (h *Handlers) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	projects, err := h.services.User.GetUserProjects(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			h.handleNotFoundError(w, r, "User")
			return
		}
		h.handleInternalError(w, r, err, "Failed to get user projects")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, projects, "User projects retrieved successfully")
}

// ListUsers handles GET /api/v1/users (with pagination)
func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.getPaginationParams(r)

	users, err := h.services.User.ListUsers(r.Context(), limit, offset)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to list users")
		return
	}

	response := map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"count":  len(users),
	}

	h.writeSuccessResponse(w, r, http.StatusOK, response, "Users retrieved successfully")
}
