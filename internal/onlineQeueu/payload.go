package onlineQeueu

type JoinQueueRequest struct {
	OfficeID     uint `json:"office_id"`
	ClientNumber int  `json:"client_number"`
}

type CancelQueueRequest struct {
	OfficeID     uint `json:"office_id"`
	ClientNumber int  `json:"client_number"`
}

type TakeClientRequest struct {
	OfficeID   uint `json:"office_id"`
	OperatorID uint `json:"operator_id"`
}

type FinishServiceRequest struct {
	OfficeID   uint `json:"office_id"`
	OperatorID uint `json:"operator_id"`
}
