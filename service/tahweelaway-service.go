package service

import (
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
)

type TahweelawayService struct {
	persistence persistence.PersistenceLayerInterface
}

// NewService returns new advisor manager that allows CRUD operations
func NewService(persistence persistence.PersistenceLayerInterface) *TahweelawayService {
	return &TahweelawayService{
		persistence: persistence,
	}
}
