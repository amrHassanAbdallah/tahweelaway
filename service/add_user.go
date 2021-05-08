package service

import (
	"context"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
)

func (s TahweelawayService) AddUser(ctx context.Context, user User) (*persistence.User, error) {
	return s.persistence.AddUser(ctx, persistence.User{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Currency: user.Currency,
	})
}
