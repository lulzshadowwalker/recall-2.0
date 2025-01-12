package handler

import (
	"net/http"
)

type HomeHandler struct {
	// service service
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Recall HTTP interface"))
}
