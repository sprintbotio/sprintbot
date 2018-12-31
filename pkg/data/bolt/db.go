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
	teamsBucket = []byte("teams")
	usersBucket  = []byte("users")
	standupSchedule = []byte("standup_schedule")
	standupLogs = []byte("standup_logs")
	buckets = [][]byte{
		usersBucket,
		teamsBucket,
		standupSchedule,
		standupLogs,
	}
)

func Setup()error  {
	if db == nil{
		return errors.New("no db connection could not complete setup")
	}
	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets{
			if _, err := tx.CreateBucketIfNotExists(b); err != nil{
				return errors.Wrap(err, "failed to create "+string(b)+" bucket")
			}
		}

		return nil
	})
}

