# resterror
REST API error handling in Go.

## Motivation
There is a tension between making errors helpful and concise for the end user versus making them expansive and analytic for the operator. Too often the operator wins and the errors are overly verbose, to the point of including stack traces or other overwhelming detail.

This package is an attempt to service each consumer role. The reported errors are resonable concise, concentrating on information the user should find helpful. But they also contain internal details such as method names an operator might find diagnostic but not in a way that overwhelms the user. 

Stack trace errors are not concise and don't concentrate on information that the consumer finds helpful. That is why this package is designed using error nesting that behaves as an operational stack trace, showing the path through the elements of the system rather than as an execution trace, showing the path through the code.

There is a lack of finding standardized error messaging and handling for REST APIs, especially using Golang. Golang's error handling is core to the language, but don't describe how to handle errors. There have been community discussions on improving standardization of error handling, but many miss the errors within the application domain. That is, errors are important as your types.

Error handling is critical to application design and is complicated by the various consumer roles that require error information. By considering error codes, messages, and logical stack traces in our design, we can fulfill the needs of each consumer. By integrating our Error into our domain, we give all parts of our application a common language to communicate about when unexpected things happen.

## Notes
Errors serve different goals for each of its consumer roles -- the application, the end user, and the operator. This application descibes the **INSERT ROLE HERE**. Errors can be divided into two categories: well-defined, and undefined errors. Well-defined errors are specified by the API, such as os.IsNotExist() error returned from os.Open(). These allow us to manage our application flow because we know what to expect and can work with them on a case-by-case basis. Undefined errors are undocumented by the API and therefore we are unable to thoughtfully handle it. This can occur from poor documentation, or can occur when APIs we depend on add additional error conditions after we've integrated our code with them. **This is why I tried to use the least amount of third-party software in this package, if any.**

Errors need to be different things to different consumers of them. In a given system, we have at least 3 consumer roles: application, end-user, and the operator.

An unexpected benefit to custom error handling is the ease of writing error-dependent tests, as well as write error-sensitive code outside of tests. Error.Is enables ease of testing by returning a boolean reporting whether the argument is of type *Error, and if so, its Code field has the specified value.

Is function makes it straightforward for code to change behaviour depending on the error condition, such as in the face of a permission error as opposed to a network error:

if Error.Is(errors.EEXIST, err) { ... }


### Consumer roles

#### Application
The first line of defense in error handling is the application itself. Your application code can recover from error states quickly and without paging anyone in the middle of the night. However, application error handling is the least flexible and it can only handle well-defined error states. An example of this is your web browser receiving a 301 redirect code and navigating you to a new location. It's seamless and most people are oblivious to it. It's able to do this because the HTTP spec has well-defined error codes.

#### End User
If your application can't handle the error, hopefully your end user can resolve the issue. Your end user can see an error state such as "Your debit card is declined" and is flexible enough to resolve it (i.e. deposit money into their bank account).

Unlike the application role, the end user needs a human-readable message that can provide context to help them resolve the error.

These users are still limited to well-defined errors since revealing undefined errors could compromise the security of your system. For example, a postgres error may detail query or schema information that can be used by an attacker. When confronted with an undefined error, it may be appropriate to simple tell the user to contact trechnical support.

#### Operator
The last line of defense is the system operator which may be a developer or ops person. These people understand the system and can work with any kind of error. In this role, you want to see as much information as possible. In addition to the error code and human-readable message, a logical stack trace can help the operator understand the program flow.

## Translating error codes to our domain
The error codes are defined in codes.go. Instead of creating more fine-grained error codes, it's easier to manage more generic codes. We use HTTP Status Code Definitions.

Error codes are specific to this application, so when we interact with external libraries we must translate those errors to our domain's error codes. For example, if our application implements our UserService in Postgres, we will need to translate a sql.ErrNoRows and it would break down if we also implement UserService with a non-SQL database. There is an example of this in serviceExample.go.

## Tech/framework used
TBD

## TODO
1) Use https://github.com/pkg/errors to add stack trace and for wrapping errors.

2) Define an interface for internal/server errors.

3) Make error messages more descriptive and keep them as constants in a file.

4) Use an error tracking platform such as Sentry for logging.

5) Use an audit logger.

6) When a model has an error and includes Op, the HTTP request should not show this. There needs to be a way to ensure clients do not see these internal errors.

## Resources
https://medium.com/@ozdemir.zynl/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd

https://middlemost.com/failure-is-your-domain/

https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

https://github.com/pkg/errors

