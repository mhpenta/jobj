package safeunmarshal

import "errors"

// ErrExpectedJSONArray is returned when the response being unmarshalled is not a JSON array in our unmarshaller but is
// an array type
var ErrExpectedJSONArray = errors.New("expected JSON array for array type")

// ErrJSONRepairFailed is returned when JSON repair attempts fail
var ErrJSONRepairFailed = errors.New("JSON repair failed")
