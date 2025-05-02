package onlineQeueu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"onlineQueue/pkg/req"
	"onlineQueue/pkg/res"
	"strconv"
)

type QueueHandler struct {
	QueueService *QueueService
}

type QueueHandlerDeps struct {
	QueueService *QueueService
}

func NewQueueHandler(router *http.ServeMux, deps QueueHandlerDeps) {
	handler := &QueueHandler{
		QueueService: deps.QueueService,
	}

	router.HandleFunc("POST /queue/join", handler.JoinQueue())
	router.HandleFunc("POST /queue/cancel", handler.CancelQueue())
	router.HandleFunc("POST /queue/take", handler.TakeClient())
	router.HandleFunc("POST /queue/finish", handler.FinishService())
	router.HandleFunc("GET /queue/position", handler.GetClientPosition())
}

// JoinQueue добавляет клиента в очередь по заданному офису.
//
// @Summary      Присоединиться к очереди
// @Description  Добавляет клиента с указанным номером в очередь определенного офиса
// @Tags         Очередь
// @Accept       json
// @Produce      json
// @Param        request  body      JoinQueueRequest  true  "Данные клиента и офиса"
// @Success      201      {object}  map[string]string  "Клиент добавлен"
// @Failure      400      {string}  string             "Неверный формат запроса"
// @Failure      500      {string}  string             "Ошибка сервера или бизнес-логики"
// @Router       /queue/join [post]
func (h *QueueHandler) JoinQueue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[JoinQueueRequest](&w, r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.QueueService.AddClientToQueue(body.OfficeID, body.ClientNumber); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]string{"message": "Client added to queue"}, http.StatusCreated)
	}
}

func (h *QueueHandler) CancelQueue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CancelQueueRequest](&w, r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.QueueService.RemoveClientFromQueue(body.OfficeID, body.ClientNumber); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]string{"message": "Client removed from queue"}, http.StatusOK)
	}
}

func (h *QueueHandler) TakeClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[TakeClientRequest](&w, r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		clientNumber, err := h.QueueService.MoveClientToInService(body.OfficeID, body.OperatorID)
		if err != nil {
			http.Error(w, "Queue is empty or error occurred", http.StatusBadRequest)
			return
		}

		res.Json(w, map[string]interface{}{
			"message":       "Client taken for service",
			"client_number": clientNumber,
		}, http.StatusOK)
	}
}

func (h *QueueHandler) FinishService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[FinishServiceRequest](&w, r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.QueueService.FinishService(body.OfficeID, body.OperatorID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]string{"message": "Service finished"}, http.StatusOK)
	}
}

// GetClientPosition возвращает позицию клиента в очереди.
//
// @Summary      Узнать место клиента в очереди
// @Tags         Очередь
// @Accept       json
// @Produce      json
// @Param        office_id    query     int    true  "ID офиса"
// @Param        phone        query     string true  "Номер телефона клиента"
// @Success      200 {object} GetQueueResponse "Информация о позиции клиента"
// @Failure      400 {string} string "Некорректные параметры запроса"
// @Router       /queue/position [get]
func (h *QueueHandler) GetClientPosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметры из URL
		officeIDStr := r.URL.Query().Get("office_id")
		phoneNumber := r.URL.Query().Get("phone")

		if officeIDStr == "" || phoneNumber == "" {
			http.Error(w, "Missing office_id or phone", http.StatusBadRequest)
			return
		}

		officeID, err := strconv.Atoi(officeIDStr)
		if err != nil {
			http.Error(w, "Invalid office_id format", http.StatusBadRequest)
			return
		}

		keyQueue := fmt.Sprintf("queue:%d", officeID)
		queue, err := h.QueueService.RedisClient.LRange(context.Background(), keyQueue, 0, -1).Result()
		if err != nil {
			http.Error(w, "Failed to get queue", http.StatusInternalServerError)
			return
		}

		clientsBefore := -1
		for index, number := range queue {
			if number == phoneNumber {
				clientsBefore = index + 1
				break
			}
		}

		response := map[string]interface{}{
			"общее колво": len(queue),
		}

		if clientsBefore == -1 {
			response["сообщение"] = "Вас нет в очереди"
		} else {
			response["лично ваша очередь"] = clientsBefore
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}
}
