package offices

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/internal/onlineQeueu"
	"onlineQueue/pkg/req"
	"onlineQueue/pkg/res"
)

type OfficeHandler struct {
	OfficeRepository *OfficeRepository
	QueueService     *onlineQeueu.QueueService
}

type OfficeHandlerDeps struct {
	OfficeRepository *OfficeRepository
	Config           *configs.Config
	QueueService     *onlineQeueu.QueueService
}

func NewOfficeHandler(router *http.ServeMux, deps OfficeHandlerDeps) {
	handler := &OfficeHandler{
		OfficeRepository: deps.OfficeRepository,
		QueueService:     deps.QueueService,
	}
	router.HandleFunc("POST /register", handler.CreateOffice())
}

func (handler *OfficeHandler) CreateOffice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[OfficeCreateRequest](&w, r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		office := NewOffice(body.Name, body.Address, body.WorkingHours)
		createdOffice, err := handler.OfficeRepository.CreateOffice(office)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.QueueService.CreateOfficeQueue(createdOffice.ID)
		if err != nil {
			http.Error(w, "Failed to create queue", http.StatusInternalServerError)
			return
		}
		res.Json(w, createdOffice, http.StatusCreated)

	}
}
