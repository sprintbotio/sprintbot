package bolt

import (
	"encoding/json"

	"github.com/pkg/errors"
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
	if s.ID == "" {
		return errors.New("no id specified for stand up")
	}
	return sr.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(s)
		if err != nil {
			return errors.Wrap(err, " failed to marshal stand up")
		}
		return tx.Bucket(standups).Put([]byte(s.ID), data)
	})
}

func (sr *StandUpRepo) Get(sid string) (*domain.StandUp, error) {
	var su domain.StandUp
	err := sr.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(standups).Get([]byte(sid))
		if nil == data {
			return domain.NewNotFoundErr("stand up")
		}
		if err := json.Unmarshal(data, &su); err != nil {
			return errors.Wrap(err, " failed to unmarshal stand up")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &su, nil
}
