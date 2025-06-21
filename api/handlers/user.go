package handlers

import (
	"encoding/json"
	"net/http"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"github.com/nurfianqodar/school-microservices/utils/httperr"
)

type userHandler struct {
	s pbusers.UserServiceClient
}

func (h *userHandler) RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users/{$}", h.handleCreateOneUser)
}

func NewUserHandler(s pbusers.UserServiceClient) Handler {
	return &userHandler{s: s}
}

func (h *userHandler) handleCreateOneUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body := new(pbusers.CreateOneUserRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		httperr.ErrInvalidRequestBody.Send(w)
		return
	}

	res, err := h.s.CreateOneUser(r.Context(), body)
	if err != nil {
		httperr.ConvertGRPCErrorToHTTPErr(err).Send(w)
		return
	}

	json.NewEncoder(w).Encode(res)
}
