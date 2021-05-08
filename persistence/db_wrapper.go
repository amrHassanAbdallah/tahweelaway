package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	//TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(ctx context.Context, db *sql.DB) (*SQLStore, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}, nil
}

// HashPassword encrypts user password
func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (s *SQLStore) AddUser(ctx context.Context, user User) (*User, error) {
	id := uuid.New()
	err := user.HashPassword()
	if err != nil {
		return nil, err
	}
	cuser, err := s.CreateUser(ctx, CreateUserParams{
		ID:       id,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Currency: user.Currency,
	})
	if err != nil {
		if e, k := err.(*pq.Error); k {
			if e.Code == "23505" {
				err = &DuplicateEntityException{entity: "user"}
			}
		}
		/*	fmt.Printf("Type of err is %T \n", err)
			fmt.Printf("Type of err is %#v \n", err)*/
		return nil, err
	}
	return &cuser, nil
}
func (s *SQLStore) AddBank(ctx context.Context, bank CreateBankParams) (*Bank, error) {
	id := uuid.New()
	bank.ID = id

	cuser, err := s.CreateBank(ctx, bank)
	if err != nil {
		if e, k := err.(*pq.Error); k {
			if e.Code == "23505" {
				err = &DuplicateEntityException{entity: "bank"}
			}
		}
		/*	fmt.Printf("Type of err is %T \n", err)
			fmt.Printf("Type of err is %#v \n", err)*/
		return nil, err
	}
	return &cuser, nil
}
func (s *SQLStore) GetUser(ctx context.Context, ID string) (*User, error) {
	uid, err := uuid.Parse(ID)
	if err != nil {
		return nil, &InvalidIDException{}
	}
	cuser, err := s.GetUserByID(ctx, uid)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &EntityNotFound{
				name: "user",
				id:   ID,
			}
		}

		fmt.Printf("Type of err is %T \n", err)
		fmt.Printf("Type of err is %#v \n", err)
		return nil, err
	}
	return &cuser, nil
}
