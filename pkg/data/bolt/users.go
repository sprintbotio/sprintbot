package bolt

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"go.etcd.io/bbolt"
)

type UserRepository struct {
	db *bolt.DB
}

func NewUserRepository(db *bolt.DB)*UserRepository  {
	return &UserRepository{
		db:db,
	}
}

func (ur *UserRepository)AddUser(u *domain.User)(string, error){
	 err := ur.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		data, err := json.Marshal(u)
		if err != nil{
			return errors.Wrap(err, "failed to parse user")
		}
		return b.Put([]byte(u.ID),data)
	 })
	return u.ID, err
}

func (ur *UserRepository)GetUser(id string)(*domain.User, error )  {
	var u domain.User
	err := ur.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(usersBucket).Get([]byte(id))
		if len(data) == 0{
			return errors.New("failed to find user with id " + id)
		}
		if err := json.Unmarshal(data, &u); err != nil{
			return err
		}
		return nil
	})
	return &u, err
}

func (ur *UserRepository)Update(u *domain.User)error  {
	if u.ID == ""{
		return errors.New("cannot update user with no id")
	}
	err := ur.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(u)
		if err != nil{
			return errors.Wrap(err, "failed to parse user")
		}
		return tx.Bucket(usersBucket).Put([]byte(u.ID),data)
	})
	return err
}
