package bolt

import (
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"go.etcd.io/bbolt"
)

type StandUpRepo struct {
	db *bolt.DB
}

func NewStandUpRepo(db *bolt.DB) *StandUpRepo {
	return &StandUpRepo{db: db}
}

func (sr *StandUpRepo) SaveUpdate(s *domain.StandUp) error {
	return nil
}
