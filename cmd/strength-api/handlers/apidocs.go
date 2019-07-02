package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savaki/swag/swagger"

	"github.com/savaki/swag"

	"github.com/savaki/swag/endpoint"
)

func CreateSwaggerAPIDocs(b *Book, router *gin.Engine) *swagger.API {
	remaining := endpoint.New("get", "/remaining", "Get Remaining Seat(s)",
		endpoint.Handler(b.Remaining),
		endpoint.Response(http.StatusOK, RemainingSeats{}, "Remaining Seats"),
	)
	book := endpoint.New("post", "/book", "Reserve a seat",
		endpoint.Handler(b.Book),
		endpoint.Body(RequestSeat{}, "A seat to reserve", true),
		endpoint.Response(http.StatusOK, ReservedSeat{}, "Reserved Seat"),
	)

	apiDocs := swag.New(
		swag.Endpoints(remaining, book),
	)
	apiDocs.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
	})
	return apiDocs
}
