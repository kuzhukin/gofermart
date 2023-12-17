package client

import "errors"

const (
	RegistredStatus  = "REGISTERED"
	InvalidStatus    = "INVALID"
	ProcessingStatus = "PROCESSING"
	ProcessedStatus  = "PROCESSED"
)

const (
	succesHandlingStatusCode      = 200
	orderIsNotRegistredStatusCode = 204
	requestLimitExceded           = 429
)

var ErrOrderIsNotRegistredInAccural = errors.New("order isn't registred in accural system")
var ErrRequestsLimitExceeded = errors.New("accural system request limit exceeded")

type AccuralResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accural int    `json:"accural"`
}
