package safeunmarshall

import "errors"

// ErrExpectedJSONArray is returned when the response being unmarshalled is not a JSON array in our unmarshaller
// utilities
var ErrExpectedJSONArray = errors.New("expected JSON array for array type")
