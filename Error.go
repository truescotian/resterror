// Baseline Application Error Type.
//
// Given an understading that we need error codes, human-readable messages, and
// logical stack trace, we can construct a simple error type to handle most of our
// application's errors.

package main

// Error is the center of this package and is a concrete representation of our errors.
// it has many fields, any of which can be left unset.
//
// Implements ClientError interface and is the baseline error type for this application.
//
// Code and Message fields provide communication to our application and end user
// roles. Op and Err fields allow us to chain errors together so that we can build
// the logical stack trace for our operator.
type Error struct {
	// Machine-readable error code. This is similar to the "kind" of error
	// such as ENOTFOUND, EEXISTS.
	Code string

	// Human-readable message.
	// This could be returned in JSON response. Ex: { "detail": "Wrong password" }.
	Message string

	// Op is a logical operation. It denotes the operation being performed.
	// Typically holds the name of the method or function reporting the error.
	Op string

	// Err is the original error (unmarshall errors, network errors...) which
	// caused this error, set it to nil if there isn't any.
	Err error
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

// NewError returns an Error using the passed arguments.
func NewError(op string, message string, code string, err error) Error {
	return &Error{
		Op:      op,
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// ErrorCode returns the code of the root error if available. Otherwise returns EINTERNAL.
//
// This function allows for working with Error effectively, avoiding issues such as type asserting
// whenever we want to access Error.Code. This and other issues are solved by the following:
//
// 1. Return no error code for nil errors.
// 2. Search the chain of Error.Err until a defined Code is found.
// 3. If no code is defined then return an internal error code (EINTERNAL).
func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}

	return EINTERNAL
}

// ErrorMessage is a utility function to extract error messages from error
// values.
//
// This is similar to ErrorCode except for the following rules:
//
// 1. Returns no error message for nil errors.
// 2. Searches the chain of Error.Err until a defined Message is found.
// 3. If no message is defined then return a generic error message.
//
// Returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

// Error method is used to return an error string suitable for operators.
// There's no definitive standard for how to format this message, but
// these are formatted here with these goals in mind:
//
// 1. Show the logical stack trace first. It provides context for the rest
// of the message. It also allows us to sort error lines to group them together.
// 2. Show Code, Message at the end.
// 3. Print on a single line so it's easy to grep.
//
// NOTE(truescotian): This implementation assumes that Err cannot coexist with Code
// or Message on any given error.
//
// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
//
// Source: https://upspin.googlesource.com/upspin/+/033a63d02f07/errors/errors.go#484
func Is(code string, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Code != OTHER {
		return e.Code == code
	}
	if e.Err != nil {
		return Is(code, e.Err)
	}
	return false
}
