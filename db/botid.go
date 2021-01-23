package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("botid")
var db *bolt.DB

type BotID struct {
	Key      int
	IdString string
}

type ChannelInfo struct {
	Key       string
	ChannelID string
}

func Init(dbpath string) error {
	var err error
	db, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func SetBotID(botId string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Put([]byte("botid"), []byte(botId))
	})
}

func GetBotID(kv string) (string, error) {
	var bottoken string
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == kv {
				bottoken = string(v)
				break
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return bottoken, nil
}

func GetChannels() (map[string]string, error) {
	mapChannel := make(map[string]string)
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) != "botid" {
				mapChannel[string(k)] = string(v)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mapChannel, nil
}

func AddChannel(name string, id string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Put([]byte(name), []byte(id))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
