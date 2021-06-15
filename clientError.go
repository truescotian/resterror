package main

// ClientError is one of the two main error types (Client Error for 4xx and
// Server Error for 5xx). Here we can declare interfaces based on the behaviour
// we expect from these two types and use type assertion on rootHandler
// to make decisions about the error.
//
// This is a strong definition for errors so it's easy to define an interface
// based on the behaviour we expect from the error and assert for this interface
// in the main handler.
type ClientError interface {
	Error() string
	// ResponseBody returns response body of the error (title, message, error code...)
	// in bytes.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code (4xx, 5xx) and headers
	// (content-type, no-cache).
	ResponseHeaders() (int, map[string]string)
}
