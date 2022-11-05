package errors

type ResponseError struct {
	Reason string `json:"reason"`
}