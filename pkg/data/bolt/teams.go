package bolt

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"go.etcd.io/bbolt"
)

type TeamRepository struct {
	db *bolt.DB
}

func NewTeamRespository(db *bolt.DB)*TeamRepository  {
	return &TeamRepository{
		db:db,
	}
}

func (tr *TeamRepository)AddTeam(team domain.Team)(string,error)  {
	err := tr.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(teamsBucket)
		data, err := json.Marshal(team)
		if err != nil{
			return errors.Wrap(err, "failed to unmarshal team")
		}
		return b.Put([]byte(team.ID),data)
	})
	return team.ID, err
}

func (tr *TeamRepository)GetTeam(id string)(*domain.Team, error )  {
	var t = domain.Team{}
	tr.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(teamsBucket).Get([]byte(id))
		if len(data) == 0{
			return errors.New("no such team")
		}
		if err := json.Unmarshal(data,t); err != nil{
			return errors.Wrap(err, "failed to decode team")
		}
		return nil
	})
	return &t,nil
}