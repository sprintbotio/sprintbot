package bolt

import (
	"encoding/json"
	"fmt"
	"time"

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

func (sr *StandUpRepo) GenerateID(teamID string, t time.Time) string {
	return fmt.Sprintf("%s-%d-%02d-%02d", teamID, t.Year(), t.Month(), t.Day())
}

func (sr *StandUpRepo) Delete(id string) error {
	return sr.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(standups).Delete([]byte(id))
	})
}

func (sr *StandUpRepo) List(teamID string) ([]*domain.StandUp, error) {
	var standUps []*domain.StandUp
	err := sr.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(standups).ForEach(func(k, v []byte) error {
			s := &domain.StandUp{}
			if err := json.Unmarshal(v, s); err != nil {
				return err
			}
			if s.TeamID == teamID {
				standUps = append(standUps, s)
			}
			return nil
		})
	})
	return standUps, err
}

func (sr *StandUpRepo) FindByTeam(tid string) (*domain.StandUp, error) {
	var s = &domain.StandUp{}
	err := sr.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(standups).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			if s.TeamID == tid {
				break
			}
		}
		return nil
	})
	return s, err
}
