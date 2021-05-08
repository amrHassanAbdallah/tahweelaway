package service

import (
	"context"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
)

func (s TahweelawayService) GetUser(ctx context.Context, ID string) (*persistence.User, error) {
	return s.persistence.GetUser(ctx, ID)
}
