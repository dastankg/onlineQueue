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
	router.HandleFunc("POST /office", handler.CreateOffice())
	router.HandleFunc("GET /offices", handler.GetOffices())
}

// CreateOffice создает новый офис
//
// @Summary 	Cоздать новый офис
// @Description Создает офис с заданным названием, адресом и графиком работы. Также создается очередь для офиса.
// @Tags 		Offices
// @Accept      json
// @Produce 	json
// @Param		body body OfficeCreateRequest  true "Данные для создание"
// @Success     201   {object}  Office                "Созданный офис"
// @Failure     400   {string}  string        "Неверный запрос или ошибка валидации"
// @Failure     500   {string}  string       "Ошибка при создании"
// @Router      /office [post]
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

// GetOffices возвращает список всех офисов.
//
// @Summary      Получить список офисов
// @Description  Возвращает массив офисов, доступных для записи в очередь
// @Tags         Offices
// @Produce      json
// @Success      200  {object}  OfficesGetResponse  "Список офисов"
// @Router       /offices [get]
func (handler *OfficeHandler) GetOffices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offices := handler.OfficeRepository.GetOffices()

		data := OfficesGetResponse{
			Offices: offices,
		}
		res.Json(w, data, http.StatusOK)
	}
}
