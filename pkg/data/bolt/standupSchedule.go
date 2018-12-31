package bolt

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"go.etcd.io/bbolt"
)

type StandUpRepository struct {
	db *bolt.DB
}

func NewStandUpRepository(db *bolt.DB) *StandUpRepository {
	return &StandUpRepository{db: db}
}

func (sr *StandUpRepository) SaveUpdate(teamID string, schedule domain.StandupSchedule) error {
	err := db.Update(func(tx *bolt.Tx) error {
		schedule.TeamID = teamID
		d, err := json.Marshal(schedule)
		if err != nil {
			return errors.Wrap(err, "failed to save standup schedule")
		}
		return tx.Bucket(standupSchedule).Put([]byte(teamID), d)
	})
	return err
}

func (sr *StandUpRepository) List() ([]*domain.StandupSchedule, error) {
	var s []*domain.StandupSchedule
	err := sr.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(standupSchedule).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ss := &domain.StandupSchedule{}
			if err := json.Unmarshal(v, ss); err != nil {
				return err
			}
			ss.TeamID = string(k)
			s = append(s, ss)
		}
		return nil
	})
	return s, err
}
