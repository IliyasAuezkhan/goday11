package http
import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"day16/internal/domain"
	"day16/pkg/logger"
)

type UserUseCase interface {
	Register(email, password string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	Update(id int64, email, password *string) error
	Delete(id int64) error
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
}

type UserHandler struct {
	service UserUseCase
	log *logger.Logger
}

func NewUserHandler(service UserUseCase, log *logger.Logger) *UserHandler {
	return &UserHandler{service: service, log: log}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode registration request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		h.log.Error("User registration failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := UserResponse{ID: user.ID, Email: user.Email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) HandleUserCRUD(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if r.Method == http.MethodPost && len(pathParts) == 1 && pathParts[0] == "users" {
		h.Register(w, r)
		return
	}

	if len(pathParts) < 2 {
		http.Error(w, "Bad request: missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getByID(w http.ResponseWriter, id int64) {
	user, err := h.service.GetByID(id)
	if err != nil {
		h.log.Error("Failed to fetch user ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := UserResponse{ID: user.ID, Email: user.Email}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, req.Email, req.Password); err != nil {
		h.log.Error("Failed to update user ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *UserHandler) delete(w http.ResponseWriter, id int64) {
	if err := h.service.Delete(id); err != nil {
		h.log.Error("Failed to delete user ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}