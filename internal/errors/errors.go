package errors

// ResponseError is a type that does NOT implement error interface
// ResponseError can be used for encoding error to JSON
type ResponseError struct {
	Reason string `json:"reason"`
}

type ErrorInvalidRequestBody struct {
	description string
}

func NewErrorInvalidRequestBody(decription string) ErrorInvalidRequestBody {
	return ErrorInvalidRequestBody{description: decription}
}

func (err ErrorInvalidRequestBody) Error() string {
	if err.description == "" {
		return "invalid request body"
	} else {
		return "invalid request body: " + err.description
	}
}
