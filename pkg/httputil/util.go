package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
)

func ParseBody[TReq any](r *http.Request) (*TReq, error) {
	var req TReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, errorutil.RequestError(err)
}

func WriteResponse(
	w http.ResponseWriter,
	result *RequestResult,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Code)
	json.NewEncoder(w).Encode(result.Data)
}
