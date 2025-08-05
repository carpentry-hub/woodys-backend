package handlers

import (
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/internal/domain"
)

// CreateProject handles POST /api/v1/projects
func (h *Handlers) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req domain.ProjectCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	project, err := h.services.Project.CreateProject(r.Context(), req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to create project")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, project, "Project created successfully")
}

// GetProject handles GET /api/v1/projects/{id}
func (h *Handlers) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	project, err := h.services.Project.GetProjectByID(r.Context(), id)
	if err != nil {
		if err.Error() == "project with id "+strconv.FormatInt(id, 10)+" not found" {
			h.handleNotFoundError(w, r, "Project")
			return
		}
		h.handleInternalError(w, r, err, "Failed to get project")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, project, "Project retrieved successfully")
}

// UpdateProject handles PUT /api/v1/projects/{id}
func (h *Handlers) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.ProjectUpdateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	project, err := h.services.Project.UpdateProject(r.Context(), id, req)
	if err != nil {
		if err.Error() == "project not found" {
			h.handleNotFoundError(w, r, "Project")
			return
		}
		h.handleInternalError(w, r, err, "Failed to update project")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, project, "Project updated successfully")
}

// DeleteProject handles DELETE /api/v1/projects/{id}
func (h *Handlers) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	err = h.services.Project.DeleteProject(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "project not found" {
			h.handleNotFoundError(w, r, "Project")
			return
		}
		h.handleInternalError(w, r, err, "Failed to delete project")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusNoContent, nil, "Project deleted successfully")
}

// SearchProjects handles GET /api/v1/projects/search
func (h *Handlers) SearchProjects(w http.ResponseWriter, r *http.Request) {
	filters := domain.ProjectSearchFilters{
		Style:       r.URL.Query().Get("style"),
		Environment: r.URL.Query().Get("environment"),
		Materials:   r.URL.Query().Get("materials"),
		Tools:       r.URL.Query().Get("tools"),
	}

	if maxTimeStr := r.URL.Query().Get("max_time_to_build"); maxTimeStr != "" {
		if maxTime, err := strconv.Atoi(maxTimeStr); err == nil {
			filters.MaxTimeToBuild = maxTime
		}
	}

	if minRatingStr := r.URL.Query().Get("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 32); err == nil {
			filters.MinRating = float32(minRating)
		}
	}

	filters.Limit, filters.Offset = h.getPaginationParams(r)

	projects, err := h.services.Project.SearchProjects(r.Context(), filters)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to search projects")
		return
	}

	response := map[string]interface{}{
		"projects": projects,
		"filters":  filters,
		"count":    len(projects),
	}

	h.writeSuccessResponse(w, r, http.StatusOK, response, "Projects retrieved successfully")
}

// CreateComment handles POST /api/v1/projects/{project_id}/comments
func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.CommentCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	req.ProjectID = projectID

	comment, err := h.services.Comment.CreateComment(r.Context(), req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to create comment")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, comment, "Comment created successfully")
}

// GetProjectComments handles GET /api/v1/projects/{project_id}/comments
func (h *Handlers) GetProjectComments(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	comments, err := h.services.Comment.GetProjectComments(r.Context(), projectID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to get project comments")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, comments, "Comments retrieved successfully")
}

// DeleteComment handles DELETE /api/v1/comments/{id}
func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	err = h.services.Comment.DeleteComment(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "comment not found" {
			h.handleNotFoundError(w, r, "Comment")
			return
		}
		h.handleInternalError(w, r, err, "Failed to delete comment")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusNoContent, nil, "Comment deleted successfully")
}

// GetCommentReplies handles GET /api/v1/comments/{id}/replies
func (h *Handlers) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	replies, err := h.services.Comment.GetCommentReplies(r.Context(), id)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to get comment replies")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, replies, "Replies retrieved successfully")
}

// CreateReply handles POST /api/v1/comments/{id}/reply
func (h *Handlers) CreateReply(w http.ResponseWriter, r *http.Request) {
	parentID, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.CommentCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	reply, err := h.services.Comment.CreateReply(r.Context(), parentID, req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to create reply")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, reply, "Reply created successfully")
}

// CreateRating handles POST /api/v1/projects/{project_id}/ratings
func (h *Handlers) CreateRating(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.RatingCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	req.ProjectID = projectID

	rating, err := h.services.Rating.CreateRating(r.Context(), req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to create rating")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, rating, "Rating created successfully")
}

// UpdateRating handles PUT /api/v1/projects/{project_id}/ratings
func (h *Handlers) UpdateRating(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.RatingUpdateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	rating, err := h.services.Rating.UpdateRating(r.Context(), userID, projectID, req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to update rating")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, rating, "Rating updated successfully")
}

// GetProjectRatings handles GET /api/v1/projects/{project_id}/ratings
func (h *Handlers) GetProjectRatings(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	ratings, err := h.services.Rating.GetProjectRatings(r.Context(), projectID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to get project ratings")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, ratings, "Ratings retrieved successfully")
}

// CreateProjectList handles POST /api/v1/project-lists
func (h *Handlers) CreateProjectList(w http.ResponseWriter, r *http.Request) {
	var req domain.ProjectListCreateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	projectList, err := h.services.ProjectList.CreateProjectList(r.Context(), req)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to create project list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, projectList, "Project list created successfully")
}

// GetProjectList handles GET /api/v1/project-lists/{id}
func (h *Handlers) GetProjectList(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	projectList, err := h.services.ProjectList.GetProjectListByID(r.Context(), id, userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to get project list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, projectList, "Project list retrieved successfully")
}

// UpdateProjectList handles PUT /api/v1/project-lists/{id}
func (h *Handlers) UpdateProjectList(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req domain.ProjectListUpdateRequest
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	projectList, err := h.services.ProjectList.UpdateProjectList(r.Context(), id, req, userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to update project list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, projectList, "Project list updated successfully")
}

// DeleteProjectList handles DELETE /api/v1/project-lists/{id}
func (h *Handlers) DeleteProjectList(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	err = h.services.ProjectList.DeleteProjectList(r.Context(), id, userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to delete project list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusNoContent, nil, "Project list deleted successfully")
}

// AddProjectToList handles POST /api/v1/project-lists/{id}/projects
func (h *Handlers) AddProjectToList(w http.ResponseWriter, r *http.Request) {
	listID, err := h.parseIDFromPath(r, "id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	var req struct {
		ProjectID int64 `json:"project_id" validate:"required"`
	}
	if err := h.validateRequestBody(r, &req); err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	err = h.services.ProjectList.AddProjectToList(r.Context(), listID, req.ProjectID, userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to add project to list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusCreated, nil, "Project added to list successfully")
}

// RemoveProjectFromList handles DELETE /api/v1/project-lists/{list_id}/projects/{project_id}
func (h *Handlers) RemoveProjectFromList(w http.ResponseWriter, r *http.Request) {
	listID, err := h.parseIDFromPath(r, "list_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	projectID, err := h.parseIDFromPath(r, "project_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	// TODO: Get user ID from authentication context
	userID := int64(1) // Placeholder

	err = h.services.ProjectList.RemoveProjectFromList(r.Context(), listID, projectID, userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to remove project from list")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusNoContent, nil, "Project removed from list successfully")
}

// GetUserProjectLists handles GET /api/v1/users/{user_id}/project-lists
func (h *Handlers) GetUserProjectLists(w http.ResponseWriter, r *http.Request) {
	userID, err := h.parseIDFromPath(r, "user_id")
	if err != nil {
		h.handleValidationError(w, r, err)
		return
	}

	projectLists, err := h.services.ProjectList.GetUserProjectLists(r.Context(), userID)
	if err != nil {
		h.handleInternalError(w, r, err, "Failed to get user project lists")
		return
	}

	h.writeSuccessResponse(w, r, http.StatusOK, projectLists, "User project lists retrieved successfully")
}
