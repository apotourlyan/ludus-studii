package httputil

import (
	"context"
	"net/http"
)

func HandlePost[TReq, TResp any](
	w http.ResponseWriter,
	r *http.Request,
	serviceFn func(context.Context, *TReq) (*TResp, error),
	codeMap map[string]int,
) {
	req, err := ParseBody[TReq](r)
	if err != nil {
		WriteResponse(w, ErrorResult(err, codeMap))
		return
	}

	resp, err := serviceFn(r.Context(), req)
	if err != nil {
		WriteResponse(w, ErrorResult(err, codeMap))
		return
	}

	WriteResponse(w, CreatedResult(resp))
}
