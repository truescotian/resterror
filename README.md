# resterror
REST API error handling in Go.

## Motivation
There is a lack of finding standardized error messaging and handling for REST APIs, especially using Golang.

## Tech/framework used
TBD

## Error codes (domains)
The error codes are defined in codes.go. Instead of creating more fine-grained error codes, it's easier to manager more generic codes. We use HTTP Status Code Definitions.

Error codes are specific to this application, so when we interact with external libraries we must translate those errors to our domain's error codes. For example, if our application implements our UserService in Postgres, we will need to translate a sql.ErrNoRows and it would break down if we also implement UserService with a non-SQL database.

## TODO
1) Use https://github.com/pkg/errors to add stack trace and for wrapping errors.

2) Define an interface for internal/server errors.

3) Make error messages more descriptive and keep them as constants in a file.

4) Use an error tracking platform such as Sentry for logging.

5) Use an audit logger.

## Resources
https://medium.com/@ozdemir.zynl/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd

https://middlemost.com/failure-is-your-domain/

https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html
