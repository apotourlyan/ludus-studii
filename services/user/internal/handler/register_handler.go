package handler

import (
	"net/http"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errcode"
	"github.com/apotourlyan/ludus-studii/pkg/httputil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register"
	rerrcode "github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errcode"
)

type registerHandler struct {
	service *register.Service
	codeMap map[string]int
}

func NewRegisterHandler(service *register.Service) Post {
	codeMap := errcode.GetBaseCodeMap()
	codeMap[rerrcode.EmailExists] = http.StatusConflict

	return &registerHandler{service, codeMap}
}

func (h *registerHandler) Execute(w http.ResponseWriter, r *http.Request) {
	httputil.HandlePost(w, r, h.service.Register, h.codeMap)
}
