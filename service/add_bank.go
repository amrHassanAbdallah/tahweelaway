package service

import (
	"context"
	"database/sql"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
)

func (s TahweelawayService) AddBank(ctx context.Context, bank Bank) (*persistence.Bank, error) {
	reference := sql.NullString{
		String: "",
		Valid:  false,
	}
	if bank.Reference != nil {
		reference.Valid = true
		reference.String = *bank.Reference
	}
	val := persistence.CreateBankParams{
		Name:              bank.Name,
		UserID:            bank.UserID,
		BranchNumber:      bank.BranchNumber,
		AccountNumber:     bank.AccountNumber,
		AccountHolderName: bank.AccountHolderName,
		Currency:          bank.Currency,
		ExpireAt:          bank.ExpireAt,
		Reference:         reference,
	}

	return s.persistence.AddBank(ctx, val)
}
