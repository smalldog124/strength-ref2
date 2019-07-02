package main

import (
	"github.com/bigbearsio/strength-ref2/internal/database"
	"github.com/bigbearsio/strength-ref2/cmd/strength-api/handlers"
  "github.com/gin-gonic/gin"
  "net/http"

  "github.com/boltdb/bolt"
  "log"
  "time"

  "encoding/json"

	"github.com/savaki/swag"
	"github.com/savaki/swag/endpoint"
	"github.com/savaki/swag/swagger"
)

const dbBucket = "Default"

const timeLimitMS = 10 * 1000

//https://medium.com/@ribice/serve-swaggerui-within-your-golang-application-5486748a5ed4
//https://github.com/savaki/swag/blob/master/examples/gin/main.go
func main() {
  configDB := database.Config{
    DBFile:"my.db",
    DBBucket:"Default",
    StartRow:'A',
    EndRow:'B',
    StartCol:0,
    EndCol:9,
  }
  db := database.InitDB(configDB)
  defer db.Close()

  router := gin.Default()

  routes := Routes{ db }
  bookHandler := handlers.Book{
    DB:db,
    DBBucket:"Default",
  }
  
  api := createSwaggerAPIDocs(&routes)
  api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
  })
  
  router.GET("/swagger.json", gin.WrapH(api.Handler(true)))
  router.POST("/book",bookHandler.Book)
  router.Static("/swagger/", "swagger-ui")

  router.Run(":3000")
}


func createSwaggerAPIDocs(r *Routes) *swagger.API {
	remaining := endpoint.New("get", "/remaining", "Get Remaining Seat(s)",
		endpoint.Handler(r.Remaining),
		endpoint.Response(http.StatusOK, RemainingSeats{}, "Remaining Seats"),
  )
  apiDocs := swag.New(
		swag.Endpoints(remaining),
  )
  
  return apiDocs
}

type RequestSeat struct {
  Seat string `json:"seat"`
}

type ReservedSeat struct {
  Success bool `json:success`
  Seat string `json:seat,omitempty"`
  ReserveExpiredTime int64 `json:reserve_expired_time,omitempty`
}

type RemainingSeats struct {
  UnconfimedTicketsCount int `json:"unconfimedTicketsCount"`
  Seats []string `json:"seats"`
}

type Routes struct {
  db *bolt.DB
}


func (r *Routes) Remaining(c *gin.Context) {
  result := RemainingSeats{ }
  now := time.Now()

  err := r.db.View(func(tx *bolt.Tx) error {
    // Assume bucket exists and has keys
    b := tx.Bucket([]byte(dbBucket))
  

    err := b.ForEach(func(k, v []byte) error {
      var seating = Seating{}
      json.Unmarshal(v, &seating)

      if seating.State(now) == Reserved {
        result.UnconfimedTicketsCount++
      }
      
      if len(result.Seats) < 10 && seating.State(now) == Free {
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


///// DB //////

type SeatingState int
const (
  Free      SeatingState = 0
  Reserved  SeatingState = 1
  Booked    SeatingState = 2
)


type Seating struct {
  ExpireTimestamp int64 `json:"expireTimestamp"`
  Booked bool `json:"booked"`
}

func (s *Seating) State(now time.Time) SeatingState {
  if s.Booked { 
    return Booked
  }

  if s.ExpireTimestamp > 0 && (getTimestamp(now) - s.ExpireTimestamp) < timeLimitMS   {
    return Reserved
  } else {
    return Free
  }
}


//// Funcs ////
func getTimestamp(d time.Time) int64 {
  return d.UnixNano() / int64(time.Millisecond)
}