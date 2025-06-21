package handlers

import "net/http"

type Handler interface {
	RegisterRouter(mux *http.ServeMux)
}
