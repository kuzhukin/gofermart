package client

import "errors"

const (
	RegistredStatus  = "REGISTERED"
	InvalidStatus    = "INVALID"
	ProcessingStatus = "PROCESSING"
	ProcessedStatus  = "PROCESSED"
)

var (
	ErrOrderIsNotRegisted   = errors.New("order isn't registred")
	ErrRequestLimitExceeded = errors.New("request limit exceeded")
	ErrUnknownStatusCode    = errors.New("unknown status code")
)

const (
	succesHandlingStatusCode      = 200
	orderIsNotRegistredStatusCode = 204
	requestLimitExcededStatusCode = 429
)

var ErrOrderIsNotRegistredInAccrual = errors.New("order isn't registred in accrual system")
var ErrRequestsLimitExceeded = errors.New("accrual system request limit exceeded")

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
