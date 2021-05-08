package persistence

import (
	"context"
	"fmt"
)

const (
	FROM_BANK_TO_ACCOUNT_DEPOSIT     = "FROM_BANK_TO_ACCOUNT_DEPOSIT"
	FROM_ACCOUNT_TO_ACCOUNT_TRANSFER = "FROM_ACCOUNT_TO_ACCOUNT_TRANSFER"
)

type PersistenceLayerInterface interface {
	AddUser(ctx context.Context, user User) (*User, error)
	AddBank(ctx context.Context, bank CreateBankParams) (*Bank, error)
	GetUser(ctx context.Context, ID string) (*User, error)
	GetUserBankByID(ctx context.Context, arg GetUserBankByIDParams) (Bank, error)
	TransferTx(ctx context.Context, arg CreateTransferParams) (*Transfer, error)
}
type DuplicateEntityException struct {
	entity string
}

func (nc *DuplicateEntityException) Error() string {
	return fmt.Sprintf("duplicate %s", nc.entity)
}

type InvalidIDException struct {
}

func (nc *InvalidIDException) Error() string {
	return "given id is not of type uuid."
}

type EntityNotFound struct {
	name string
	id   string
}

func (nc *EntityNotFound) Error() string {
	return fmt.Sprintf("no %s found with this id (%s)", nc.name, nc.id)
}
