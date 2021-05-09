package api

import (
	"github.com/amrHassanAbdallah/tahweelaway/service"
	"github.com/amrHassanAbdallah/tahweelaway/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

// ClientError is an error whose details to be shared with client.
type ClientError interface {
	Error() string
	// ResponseBody returns response body.
	ResponseBody() (interface{}, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}

// ValidationError implements ClientError interface.
type ValidationError struct {
	Cause  error       `json:"-"`
	Detail interface{} `json:"errors"`
	Status int         `json:"-"`
}

func (e *ValidationError) Error() string {
	if e.Cause == nil {
		return "unkown error"
	}
	return e.Detail.(string) + " : " + e.Cause.Error()
}

// ResponseBody returns JSON response body.
func (e *ValidationError) ResponseBody() (interface{}, error) {

	var validatinErrors validator.ValidationErrors
	if errors.As(e.Cause, &validatinErrors) {
		details := make([]string, 0)
		for _, e := range validatinErrors {
			details = append(details, e.Translate(utils.Translator))
		}
		e.Detail = details
	} else {
		e.Detail = []string{e.Cause.Error()}
	}

	return map[string]interface{}{"errors": e.Detail}, nil
}

// ResponseHeaders returns http status code and headers.
func (e *ValidationError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewHTTPError(err error, status int, detail interface{}) error {
	return &ValidationError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, errC error) {
	var clientError ClientError
	if !errors.As(errC, &clientError) {
		var serverError service.ServerError
		status := http.StatusInternalServerError
		message := "internal server error"
		if errors.As(errC, &serverError) {
			status = serverError.ErrorType()
			message = errC.Error()
		}
		// If the error is not ClientError, assume that it is ServerError.
		render.Status(r, status) // return 500 Internal Server Error.
		render.JSON(w, r, map[string][]string{
			"errors": {message},
		})
		return
	}
	body, err := clientError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	status, headers := clientError.ResponseHeaders() // Get http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	render.Status(r, status)
	render.JSON(w, r, body)
}
