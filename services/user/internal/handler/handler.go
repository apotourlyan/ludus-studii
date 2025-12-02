package handler

import "net/http"

type Post interface {
	Execute(w http.ResponseWriter, r *http.Request)
}
