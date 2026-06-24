// Package logger provides functions to create
// a uniform logger used in the whole application
//
// The logger will also handle correlation_ids loaded
// straight from the context
package logger

// Keys used in logging to maintain consistency
// across all components
const (
	KeyCorrelationID = "correlation_id"
	KeyService       = "service"
	KeyError         = "error"
	KeyMethod        = "method"
	KeyStatus        = "status"
	KeyDuration      = "duration"
)

// create a private type and use it as a key
// to avoid collisions in the context for the
// correlation id
type correlationIDKeyType struct{}

var correlationIDKey = correlationIDKeyType{}
