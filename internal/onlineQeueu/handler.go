package onlineQeueu

import (
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
	router.HandleFunc("GET /queue/current", handler.CurrentClient())
}

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

func (h *QueueHandler) CurrentClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		officeIDStr := r.URL.Query().Get("office_id")
		operatorIDStr := r.URL.Query().Get("operator_id")

		officeID, err := strconv.Atoi(officeIDStr)
		if err != nil {
			http.Error(w, "Invalid office_id", http.StatusBadRequest)
			return
		}
		operatorID, err := strconv.Atoi(operatorIDStr)
		if err != nil {
			http.Error(w, "Invalid operator_id", http.StatusBadRequest)
			return
		}

		clientNumber, err := h.QueueService.GetClientInService(uint(officeID), uint(operatorID))
		if err != nil {
			http.Error(w, "No client in service", http.StatusNotFound)
			return
		}

		res.Json(w, map[string]interface{}{
			"client_number": clientNumber,
		}, http.StatusOK)
	}
}
