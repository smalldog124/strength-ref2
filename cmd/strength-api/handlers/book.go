package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bigbearsio/strength-ref2/internal/book"

	"github.com/gin-gonic/gin"

	"github.com/boltdb/bolt"
)

type RequestSeat struct {
	Seat string `json:"seat"`
}

type ReservedSeat struct {
	Success            bool   `json:success`
	Seat               string `json:seat,omitempty"`
	ReserveExpiredTime int64  `json:reserve_expired_time,omitempty`
}

type RemainingSeats struct {
	UnconfimedTicketsCount int      `json:"unconfimedTicketsCount"`
	Seats                  []string `json:"seats"`
}

type Book struct {
	DB       *bolt.DB
	DBBucket string
}

func (b *Book) Book(c *gin.Context) {
	var request RequestSeat
	var result ReservedSeat
	now := time.Now()

	c.BindJSON(&request)
	log.Println("body request book", request)
	err := b.DB.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(b.DBBucket))
		v := b.Get([]byte(request.Seat))

		if v == nil {
			result = ReservedSeat{false, "", 0}
			return fmt.Errorf("Book if v")
		}

		var seating book.Seating
		json.Unmarshal(v, &seating)

		if seating.State(now) == book.Free {
			newV := book.Seating{book.GetTimestamp(now) + book.TimeLimitMS, false}
			newVBytes, _ := json.Marshal(newV)
			b.Put([]byte(request.Seat), newVBytes)

			result = ReservedSeat{true, request.Seat, newV.ExpireTimestamp}
		} else {
			result = ReservedSeat{false, "", 0}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)
}

func (b *Book) Remaining(c *gin.Context) {
	result := RemainingSeats{}
	now := time.Now()

	err := b.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(b.DBBucket))

		err := b.ForEach(func(k, v []byte) error {
			var seating = book.Seating{}
			json.Unmarshal(v, &seating)

			if seating.State(now) == book.Reserved {
				result.UnconfimedTicketsCount++
			}

			if len(result.Seats) < 10 && seating.State(now) == book.Free {
				result.Seats = append(result.Seats, string(k))
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)
}
