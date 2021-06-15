package main

// RESTError implements ClientError interface and is the baseline error
// type for this application.
//
// Code and Message fields provide communication to our application and end user
// roles. Op and Err fields allow us to chain errors together so that we can build
// the logical stack trace for our operator.
type Error struct {
	// HTTP status code.
	Code int

	// Err is the original error (unmarshall errors, network errors...) which
	// caused this error, set it to nil if there isn't any.
	Err error

	// Message is a human-message to return in JSON response. Ex: { "detail": "Wrong password" }.
	Message string

	// Op is a logical operation and nested error.
	Op string
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func (e *Error) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

func (e *Error) ResponseHeaders() (int, map[string]string) {
	return e.Code, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewError(err error, code int, message string) error {
	return &RESTError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
