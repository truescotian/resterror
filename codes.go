package main

// Types of errors.
//
// The values of the error types are common between both
// clients and servers. Do not reorder this list of remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	ECONFLICT  = "conflict"       // Action cannot be performed
	PERMISSION = "permission"     // Permission denied.
	EINTERNAL  = "internal"       // Internal error
	EINVALID   = "invalid"        // Validation failed
	ENOTFOUND  = "not_found"      // Item does not exist
	EEXIST     = "already_exists" // Item already exists
	OTHER      = "other"          // Unclassified error.
)
