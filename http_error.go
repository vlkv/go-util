package util

import (
	"fmt"
)

type HttpError struct {
	Code int
	Message string
}

var _ error = &HttpError{}

func CreateHttpError(code int, format string, a ...interface{}) HttpError {
	this := new(HttpError)
	this.Code = code
	this.Message = fmt.Sprintf(format, a...)
	return *this
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("{Code=%d, Message=%s}", err.Code, err.Message)
}

func (err *HttpError) StatusCode() int {
	return err.Code
}

func (err *HttpError) Response() []byte {
	jsonObj := map[string]interface{}{"code": err.Code, "message": err.Message}
	return JsonEncode(jsonObj)
}

