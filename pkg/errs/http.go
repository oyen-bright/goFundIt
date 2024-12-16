package errs

import "net/http"

func NotFound(message string) Error {
	return &appError{message: message, code: http.StatusNotFound}
}

func BadRequest(message string, errors interface{}) Error {
	return &appError{message: message, code: http.StatusBadRequest}
}

func Unauthorized(err Error) Error {
	return &appError{message: "Unauthorized", err: err, code: http.StatusUnauthorized}
}

func Forbidden(message string) Error {
	return &appError{message: message, code: http.StatusForbidden}
}

func InternalServerError(err error) Error {
	return &appError{message: "Internal server error", err: err, code: http.StatusInternalServerError}
}
