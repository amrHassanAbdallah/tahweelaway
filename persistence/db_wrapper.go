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

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SQLStore) TransferTx(ctx context.Context, arg CreateTransferParams) (*Transfer, error) {
	//todo pass the money currency
	var transfer Transfer
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//todo switch on the type
		id := uuid.New()
		arg.ID = id
		transfer, err = q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}
		switch arg.Type {
		case FROM_BANK_TO_ACCOUNT_DEPOSIT:
			_, err = q.AddUserBalance(ctx, AddUserBalanceParams{
				ID: arg.ToID,
				Amount: sql.NullInt64{
					Int64: arg.Amount * 100, //POUND to ERSH
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		case FROM_ACCOUNT_TO_ACCOUNT_TRANSFER:
			_, err = q.DeductUserBalance(ctx, DeductUserBalanceParams{
				ID: arg.FromID,
				Amount: sql.NullInt64{
					Int64: arg.Amount * 100, //POUND to ERSH
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
			_, err = q.AddUserBalance(ctx, AddUserBalanceParams{
				ID: arg.ToID,
				Amount: sql.NullInt64{
					Int64: arg.Amount * 100, //POUND to ERSH
					Valid: true,
				},
			})
			if err != nil {
				return err
			}

		}

		return err
	})

	return &transfer, err
}
