package handler

import (
	"edot-test/service"
	"encoding/json"
	"io"
	"net/http"
)

type TextToJsonHandler struct {
	textService *service.TextToJsonService
}

func NewTextToJsonHandler(service *service.TextToJsonService) *TextToJsonHandler {
	return &TextToJsonHandler{
		textService: service,
	}
}

func (h *TextToJsonHandler) ConvertText(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.textService.ConvertTextToJson(string(body))

	json.NewEncoder(w).Encode(result)
}
