package apiutil

import (
	"fmt"
	"net/http"
)

// ErrorObjWO is JSON object returned in case of error
type ErrorObjWO struct {
	Error ErrorWO `json:"error"`
}

// ErrorWO is JSON encapsulating error report
type ErrorWO struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Items       []*ErrorItem `json:"items,omitempty"`
}

// ErrorItem is item in ErrorWO
type ErrorItem struct {
	Name        string      `json:"name"`
	Message     string      `json:"message,omitempty"`
	Description string      `json:"description,omitempty"`
	Object      interface{} `json:"object,omitempty"`
}

// NewValidationError is validation data error
func NewValidationError(items []*ErrorItem) *Error {
	return &Error{
		Code:     "INVALID-DATA",
		Message:  "Input data validation failed",
		HTTPCode: http.StatusBadRequest,
		Items:    items,
	}
}

func (e ErrorWO) Error() string {
	return fmt.Sprintf("CODE: %s,MSG: %s, DESC: %s,Items: %s", e.Code, e.Message, e.Description, e.Items)
}

// NewErrorWO creates new API error WO
func NewErrorWO(code string, message string, desc string, items []*ErrorItem) *ErrorWO {
	return &ErrorWO{
		Code:        code,
		Message:     message,
		Description: desc,
		Items:       items,
	}
}

// NewErrorItem is ErrorItem constructor
func NewErrorItem(name, message, description string) *ErrorItem {
	return &ErrorItem{
		Name:        name,
		Message:     message,
		Description: description,
	}
}

func (e *ErrorItem) String() string {
	return fmt.Sprintf("ITEM=[%s:%s:%s]", e.Name, e.Message, e.Description)
}

// ErrorCode is error code returned by API
type ErrorCode string

// Error is encapsulated error fulfilling GO error interface
type Error struct {
	Code     ErrorCode
	Message  string
	HTTPCode int
	Items    []*ErrorItem
}

// NewError creates new API error
func NewError(httpCode int, code ErrorCode, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

func (e *Error) String() string {
	return e.Error()
	//return fmt.Sprintf("ERR[%s:%s]", e.Code, e.Error())
}

// error interface implementation
func (e Error) Error() string {
	return fmt.Sprintf("CODE: %s,MSG: %s,HTTP STS: %d,Items: %s", e.Code, e.Message, e.HTTPCode, e.Items)
}

// NewIntError is internal server error
func NewInternalServerError(e error) *Error {
	return &Error{
		Code:     "INT_APP_ERROR",
		Message:  e.Error(),
		HTTPCode: http.StatusInternalServerError,
	}
}

// CreateErrorObjWO creates ErrorObjWO suitable for JSON serialization
func CreateErrorObjWO(e *Error) *ErrorObjWO {
	return &ErrorObjWO{
		Error: ErrorWO{
			Code:    string(e.Code),
			Message: e.Message,
			Items:   e.Items,
		},
	}
}

// NewIntError is internal server error
func NewIntError(e error) *Error {
	return &Error{
		Code:     "INT_APP_ERROR",
		Message:  e.Error(),
		HTTPCode: http.StatusInternalServerError,
	}
}

var (
	// ErrCannotMatchPath is returned when path cannot be matched by path registry
	ErrCannotMatchPath = NewError(http.StatusNotFound,
		"UNKNOWN_PATH", "Cannot match path")
	// ErrFailedEncoding when JSON encoding/decoding fails
	ErrFailedEncoding = NewError(http.StatusInternalServerError,
		"FAILED_ENCODING", "Failed encoding")
	// ErrBadParameter is returned when parameter has bad value
	ErrBadParameter = NewError(http.StatusNotFound,
		"BAD_PARAM", "Bad parameter")
	// ErrRequestBodyDecoding is error returned when data in request body can not be decoded
	ErrRequestBodyDecoding = NewError(http.StatusBadRequest,
		"BAD_REQUEST_BODY_FORMAT", "Unable to decode request body")
	//ErrBadPageNumber is error returned when page_number is less than 1
	ErrBadPageNumber = NewError(http.StatusBadRequest,
		"BAD_PARAM_PAGE_NUMBER", "Bad parameter 'page_number': value must be greather than 0")
	//ErrBadPageSize is error returned when page_size is less than 1
	ErrBadPageSize = NewError(http.StatusBadRequest,
		"BAD_PARAM_PAGE_SIZE", "Bad parameter 'page_size': value : must be greather than 0")
)
