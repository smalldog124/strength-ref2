package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bigbearsio/strength-ref2/internal/book"

	"github.com/boltdb/bolt"
)

type Config struct {
	DBFile   string
	DBBucket string
	StartRow rune
	EndRow   rune
	StartCol int
	EndCol   int
}

func InitDB(cfg Config) *bolt.DB {
	// Delete and Re-create database
	fileErr := os.Remove(cfg.DBFile)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	db, err := bolt.Open(cfg.DBFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Populate Seats
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(cfg.DBBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for r := cfg.StartRow; r <= cfg.EndRow; r++ {
			for c := cfg.StartCol; c <= cfg.EndCol; c++ {
				key := []byte(fmt.Sprintf("%s%d", string(r), c))
				value, _ := json.Marshal(book.Seating{0, false})

				b.Put(key, value)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}
