package offices

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/internal/onlineQeueu"
	"onlineQueue/pkg/req"
	"onlineQueue/pkg/res"
	"strconv"
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
	router.HandleFunc("DELETE /office/{id}", handler.DeleteOffice())
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

// DeleteOffice удаляет офис по ID.
//
// @Summary      Удалить офис
// @Description  Удаляет офис по заданному идентификатору
// @Tags         Offices
// @Param        id   path      int  true  "ID офиса"
// @Produce      json
// @Success      200  {object}  OfficeDeleteResponse  "Офис успешно удалён"
// @Failure      400  {string}  string  "Некорректный ID офиса"
// @Failure      404  {string}  string  "Офис не найден"
// @Failure      500  {string}  string  "Внутренняя ошибка сервера"
// @Router       /office/{id} [delete]
func (handler *OfficeHandler) DeleteOffice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid office id", http.StatusBadRequest)
			return
		}
		_, err = handler.OfficeRepository.GetOfficeById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.OfficeRepository.DeleteOffice(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]string{"message": "Office deleted"}, http.StatusOK)
	}
}
