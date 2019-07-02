package main

import (
	"github.com/bigbearsio/strength-ref2/internal/database"
	"github.com/bigbearsio/strength-ref2/cmd/strength-api/handlers"
  "github.com/gin-gonic/gin"
  "net/http"
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

  bookHandler := handlers.Book{
    DB:db,
    DBBucket:"Default",
  }
  
  // api := createSwaggerAPIDocs(&bookHandler)
  // api.Walk(func(path string, endpoint *swagger.Endpoint) {
	// 	h := endpoint.Handler.(func(c *gin.Context))
	// 	path = swag.ColonPath(path)

	// 	router.Handle(endpoint.Method, path, h)
  // })
  
  // router.GET("/swagger.json", gin.WrapH(api.Handler(true)))
  router.POST("/book",bookHandler.Book)
  router.GET("/remaining",bookHandler.Remaining)
  router.Static("/swagger/", "swagger-ui")

  router.Run(":3000")
}


func createSwaggerAPIDocs(r *handlers.Book) *swagger.API {
	remaining := endpoint.New("get", "/remaining", "Get Remaining Seat(s)",
		endpoint.Handler(r.Remaining),
		endpoint.Response(http.StatusOK, RemainingSeats{}, "Remaining Seats"),
  )
  book := endpoint.New("post", "/book", "Reserve a seat",
    endpoint.Handler(r.Book),
    endpoint.Body(RequestSeat{}, "A seat to reserve", true),
		endpoint.Response(http.StatusOK, ReservedSeat{}, "Reserved Seat"),
	)

  apiDocs := swag.New(
		swag.Endpoints(remaining, book),
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

