package main

import (
	"log"

	"github.com/bigbearsio/strength-ref2/cmd/strength-api/handlers"
	"github.com/bigbearsio/strength-ref2/internal/database"
	"github.com/gin-gonic/gin"
)

const pathFileDB = "./internal/database/my.db"

//https://medium.com/@ribice/serve-swaggerui-within-your-golang-application-5486748a5ed4
//https://github.com/savaki/swag/blob/master/examples/gin/main.go
func main() {
	configDB := database.Config{
		DBFile:   pathFileDB,
		DBBucket: "Default",
		StartRow: 'A',
		EndRow:   'B',
		StartCol: 0,
		EndCol:   9,
	}
	db := database.InitDB(configDB)
	defer db.Close()

	router := gin.Default()

	bookHandler := handlers.Book{
		DB:       db,
		DBBucket: "Default",
	}

	// api := handlers.CreateSwaggerAPIDocs(&bookHandler, router)
	// router.GET("/swagger.json", gin.WrapH(api.Handler(true)))
	router.POST("/book", bookHandler.Book)
	router.GET("/remaining", bookHandler.Remaining)
	router.Static("/swagger/", "swagger-ui")

	log.Fatal(router.Run(":3000"))
}
