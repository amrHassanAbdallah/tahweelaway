package service

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

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
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=256,alphadash"`
	Username string `json:"username" validate:"required,min=5,max=256,alphadash"`
	Password string `json:"password"`
	Currency string `json:"currency" validate:"required,oneof=EGP_ERSH"`

	// timestamp full-date - RFC3339
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bank struct {
	ID                uuid.UUID    `json:"id"`
	Name              string       `json:"name" validate:"max=256,oneof=hsbc cib"`
	UserID            uuid.UUID    `json:"user_id" validate:"required"`
	BranchNumber      string       `json:"branch_number" validate:"required,min=5,max=256,alphadash"`
	AccountNumber     string       `json:"account_number" validate:"required,min=5,max=256,alphadash"`
	AccountHolderName string       `json:"account_holder_name" validate:"required,min=5,max=256,alphadash"`
	Reference         *string      `json:"reference"`
	Currency          string       `json:"currency" validate:"required,oneof=EGP"`
	ExpireAt          time.Time    `json:"expire_at" validate:"required,gt"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
}
