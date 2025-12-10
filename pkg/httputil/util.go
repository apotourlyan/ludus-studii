package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/pkg/httputil/content"
	"github.com/apotourlyan/ludus-studii/pkg/httputil/header"
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
	w.Header().Set(header.ContentType, content.ApplicationJson)
	w.WriteHeader(result.Code)
	json.NewEncoder(w).Encode(result.Data)
}
