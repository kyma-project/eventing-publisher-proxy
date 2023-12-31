package legacy

// Allowed patterns for the Event components.
const (
	AllowedEventTypeVersionChars = `[a-zA-Z0-9]+`
	AllowedEventIDChars          = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`
)

// Error messages.
const (
	ErrorMessageBadPayload          = "Bad payload syntax"
	ErrorMessageRequestBodyTooLarge = "Request body too large"
	ErrorMessageMissingField        = "Missing field"
	ErrorMessageInvalidField        = "Invalid field"
)

// Error type definitions.
const (
	ErrorTypeBadPayload          = "bad_payload_syntax"
	ErrorTypeRequestBodyTooLarge = "request_body_too_large"
	ErrorTypeMissingField        = "missing_field"
	ErrorTypeValidationViolation = "validation_violation"
	ErrorTypeInvalidField        = "invalid_field"
)

// Field definitions.
const (
	FieldEventID          = "event-id"
	FieldEventTime        = "event-time"
	FieldEventType        = "event-type"
	FieldEventTypeVersion = "event-type-version"
	FieldData             = "data"
)
