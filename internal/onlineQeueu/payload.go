package onlineQeueu

type JoinQueueRequest struct {
	OfficeID     uint   `json:"office_id"`
	ClientNumber string `json:"client_number"`
}

type CancelQueueRequest struct {
	OfficeID     uint   `json:"office_id"`
	ClientNumber string `json:"client_number"`
}

type TakeClientRequest struct {
	OfficeID   uint `json:"office_id"`
	OperatorID uint `json:"operator_id"`
}

type FinishServiceRequest struct {
	OfficeID   uint `json:"office_id"`
	OperatorID uint `json:"operator_id"`
}

type GetQueueRequest struct {
	OfficeID     uint   `json:"office_id" example:"1"`
	ClientNumber string `json:"client_number" example:"996501234567"`
}

type GetQueueResponse struct {
	OfficeID      uint   `json:"office_id"`
	ClientNumber  string `json:"client_number"`
	ClientsBefore int    `json:"clients_before"`
}
