package bolt

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func Connect(dbFile string )(*bolt.DB,error){
	var err error
	db, err =bolt.Open(dbFile,0666,nil)
	if err != nil{
		return nil,errors.Wrap(err,"failed to open db")
	}
	return db, nil
}

func Disconnect()  {
	if db != nil{
		if err := db.Close(); err != nil{
			logrus.Error("failed to close db", err)
		}
	}
}

var(
	usersBucket = []byte("users")
	teamsBucket = []byte("teams")
	timeZoneBucket = []byte("timezones")
)

func Setup()error  {
	if db == nil{
		return errors.New("no db connection could not complete setup")
	}
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(usersBucket); err != nil{
			return errors.Wrap(err, "failed to create user bucket")
		}
		if _, err := tx.CreateBucketIfNotExists(teamsBucket); err != nil{
			return errors.Wrap(err, "failed to create teams bucket")
		}
		return nil
	})
}

