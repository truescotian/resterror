// Baseline Application Error Type.
//
// Given an understading that we need error codes, human-readable messages, and
// logical stack trace, we can construct a simple error type to handle most of our
// application's errors.

package error

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Error is the center of this package and is a concrete representation of our errors.
// it has many fields, any of which can be left unset.
//
// Implements ClientError interface and is the baseline error type for this application.
//
// Kind and Message fields provide communication to our application and end user
// roles. Op and Err fields allow us to chain errors together so that we can build
// the logical stack trace for our operator.
type Error struct {
	// Machine-readable error code.
	// Example: ENOTFOUND, EEXISTS.
	Kind string

	// HTTP status code.
	Status int

	// Error message.
	// This could be human-readable, or a JSON response. Ex: { "detail": "Wrong password" }.
	Message string

	// Op is a logical operation. It denotes the operation being performed.
	// Typically holds the name of the method or function reporting the error.
	Op string

	// Err is the original error (unmarshall errors, network errors...) which
	// caused this error, set it to nil if there isn't any.
	Err error
}

func (e *Error) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

func (e *Error) ResponseHeaders() (string, map[string]string) {
	return e.Kind, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"X-Content-Type-Options", "nosniff",
	}
}

// Error method is used to return an error string suitable for operators.
// There's no definitive standard for how to format this message, but
// these are formatted here with these goals in mind:
//
// 1. Show the logical stack trace first. It provides context for the rest
// of the message. It also allows us to sort error lines to group them together.
// 2. Show Kind, Message at the end.
// 3. Print on a single line so it's easy to grep.
//
// NOTE(truescotian): This implementation assumes that Err cannot coexist with Kind
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
	// Otherwise print the error kind & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Kind != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Kind)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// NewError returns an Error using the passed arguments.
func NewError(op string, status int, message string, kind string, err error) *Error {
	return &Error{
		Op:      op,
		Status:  status,
		Message: message,
		Kind:    kind,
		Err:     err,
	}
}

// ErrorKind returns the kind of the root error if available. Otherwise returns EINTERNAL.
//
// This function allows for working with Error effectively, avoiding issues such as type asserting
// whenever we want to access Error.Kind. This and other issues are solved by the following:
//
// 1. Return no error kind for nil errors.
// 2. Search the chain of Error.Err until a defined Kind is found.
// 3. If no kind is defined then return an internal error kind (EINTERNAL).
func ErrorKind(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Kind != "" {
		return e.Kind
	} else if ok && e.Err != nil {
		return ErrorKind(e.Err)
	}

	return EINTERNAL
}

// ErrorMessage is a utility function to extract error messages from error
// values.
//
// This is similar to ErrorKind except for the following rules:
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

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
//
// Source: https://upspin.googlesource.com/upspin/+/033a63d02f07/errors/errors.go#484
func Is(kind string, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != OTHER {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}
