package service

import (
	"context"
	"fmt"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
	"github.com/amrHassanAbdallah/tahweelaway/proxies"
	"github.com/google/uuid"
)

func (s TahweelawayService) AddTransfer(ctx context.Context, loggedInUserID string, transfer Transfer) (*persistence.Transfer, error) {
	var bank *persistence.Transfer
	var err error
	uid, err := uuid.Parse(loggedInUserID)
	if err != nil {
		return nil, &ServiceError{
			Cause: fmt.Errorf("invlid user id"),
			Type:  400,
		}
	}
	//todo should check the transfer limit.
	switch transfer.Type {
	case persistence.FROM_BANK_TO_ACCOUNT_DEPOSIT:
		bank, err = s.handleBankToAccountTransfer(ctx, uid, transfer)
	case persistence.FROM_ACCOUNT_TO_ACCOUNT_TRANSFER:
		bank, err = s.handleAccountToAccountTransfer(ctx, uid, transfer)
	default:
		err = &ServiceError{
			Cause: fmt.Errorf("unsupported transfer type %s", transfer.Type),
			Type:  500,
		}

	}
	if err != nil {
		return nil, err
	}
	return bank, err
}

func (s TahweelawayService) handleBankToAccountTransfer(ctx context.Context, loggedInUserID uuid.UUID, transfer Transfer) (*persistence.Transfer, error) {
	bank, err := s.persistence.GetUserBankByID(ctx, persistence.GetUserBankByIDParams{
		UserID: loggedInUserID,
		ID:     transfer.FromID,
	})
	if err != nil {
		return nil, &ServiceError{
			Cause: fmt.Errorf("user id does not exist, or bank id not found"),
			Type:  404,
		}
	}
	//check if the bank have this amount and move it to our bank account
	err = proxies.ValidateAndMoveMoney(bank)
	if err != nil {
		return nil, &ServiceError{
			Cause: fmt.Errorf("no enough money in your bank account"),
			Type:  400,
		}
	}
	//start transaction to add the amount to the user balance while add the transfer record.
	return s.persistence.TransferTx(ctx, persistence.CreateTransferParams{
		FromID: transfer.FromID,
		ToID:   transfer.ToID,
		Amount: transfer.Amount,
		Type:   transfer.Type,
	})
}
func (s TahweelawayService) handleAccountToAccountTransfer(ctx context.Context, loggedInUserID uuid.UUID, transfer Transfer) (*persistence.Transfer, error) {
	user, err := s.persistence.GetUser(ctx, loggedInUserID.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &ServiceError{
			Cause: fmt.Errorf("user not found"),
			Type:  404,
		}
	}
	if user.Balance.Int64 <= transfer.Amount {
		return nil, &ServiceError{
			Cause: fmt.Errorf("no enough money in your balance to perform this transfer"),
			Type:  400,
		}
	}

	//start transaction to add the amount to the user balance while add the transfer record.
	return s.persistence.TransferTx(ctx, persistence.CreateTransferParams{
		FromID: transfer.FromID,
		ToID:   transfer.ToID,
		Amount: transfer.Amount,
		Type:   transfer.Type,
	})
}
