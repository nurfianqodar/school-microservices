package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"github.com/nurfianqodar/school-microservices/utils/httperr"
	"github.com/nurfianqodar/school-microservices/utils/httpres"
)

type userHandler struct {
	s pbusers.UserServiceClient
}

func (h *userHandler) RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users/{$}", h.handleCreateOneUser)
	mux.HandleFunc("GET /api/v1/users/{$}", h.handleListUser)
}

func NewUserHandler(s pbusers.UserServiceClient) Handler {
	return &userHandler{s: s}
}

func (h *userHandler) handleCreateOneUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	defer r.Body.Close()
	body := new(pbusers.CreateOneUserRequest)
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		httperr.ErrInvalidRequestBody.Send(w)
		return
	}

	// Invoke service and get reponse
	res, err := h.s.CreateOneUser(r.Context(), body)
	if err != nil {
		httperr.ConvertGRPCErrorToHTTPErr(err).Send(w)
		return
	}
	// Send response
	json.NewEncoder(w).Encode(res)
}

func (h *userHandler) handleListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	takeQuery := r.URL.Query().Get("take")
	skipQuery := r.URL.Query().Get("skip")

	limit, err := strconv.Atoi(takeQuery)
	if err != nil {
		httperr.New(http.StatusBadRequest, "take query must be integer").Send(w)
		return
	}

	offset, err := strconv.Atoi(skipQuery)
	if err != nil {
		httperr.New(http.StatusBadRequest, "skip query must be integer").Send(w)
		return
	}

	if limit == 0 {
		limit = 10
	}

	log.Println(limit)
	log.Println(offset)

	res, err := h.s.GetManyUser(r.Context(), &pbusers.GetManyUserRequest{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		httperr.ConvertGRPCErrorToHTTPErr(err).Send(w)
	}

	json.NewEncoder(w).Encode(httpres.New(true, res))
}
