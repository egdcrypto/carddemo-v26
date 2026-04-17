package handler

// AccountHandler defines the interface for account HTTP handlers.
// This interface allows us to mock the handler in tests or wrap it with middleware.
type AccountHandler interface {
	RegisterRoutes(r chi.Router)
}

// accountHandler implements AccountHandler.
type accountHandler struct {
	// Dependencies will be injected here (e.g., services)
}

// NewAccountHandler creates a new instance of the account handler.
func NewAccountHandler() AccountHandler {
	return &accountHandler{}
}

// RegisterRoutes sets up the routing for the account endpoints.
func (h *accountHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateAccount)
	r.Get("/{id}", h.GetAccount)
	r.Put("/{id}", h.UpdateAccount)
	r.Delete("/{id}", h.DeleteAccount)

	// UserProfile sub-routes
	r.Get("/{id}/profile", h.GetUserProfile)
	r.Put("/{id}/profile", h.UpdateUserProfile)
}

// Handler methods (to be implemented)
func (h *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {}
func (h *accountHandler) GetAccount(w http.ResponseWriter, r *http.Request)      {}
func (h *accountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request)   {}
func (h *accountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request)   {}
func (h *accountHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {}
func (h *accountHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {}
