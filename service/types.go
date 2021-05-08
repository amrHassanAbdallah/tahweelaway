package service

import "time"

type ServerError interface {
	Error() string
	ErrorType() int
}

type ServiceError struct {
	Cause error `json:"error"`
	Type  int   `json:"-"`
}

func (e *ServiceError) Error() string {
	return e.Cause.Error()
}

func (e *ServiceError) ErrorType() int {
	return e.Type
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name" validate:"required,min=5,max=256,alphadash"`
	Username string `json:"username" validate:"required,min=5,max=256,alphadash"`
	Password string `json:"password"`
	Currency string `json:"currency" validate:"required,oneof=EGP_ERSH"`

	// timestamp full-date - RFC3339
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
