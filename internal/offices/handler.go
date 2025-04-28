package offices

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/pkg/req"
	"onlineQueue/pkg/res"
)

type OfficeHandler struct {
	RegisterRepository *OfficeRepository
}

type OfficeHandlerDeps struct {
	RegisterRepository *OfficeRepository
	Config             *configs.Config
}

func NewOfficeHandler(router *http.ServeMux, deps OfficeHandlerDeps) {
	handler := &OfficeHandler{
		RegisterRepository: deps.RegisterRepository,
	}
	router.HandleFunc("POST /register", handler.CreateOffice())
}

func (handler *OfficeHandler) CreateOffice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[OfficeCreateRequest](&w, r)
		if err != nil {
			return
		}
		registers := NewOffice(body.Name, body.Address, body.WorkingHours)
		createRegister, err := handler.RegisterRepository.CreateOffice(registers)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createRegister, http.StatusCreated)

	}
}
