package main

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
)

func main() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("expiry_date", common.ExpiryDateMMYY)
		validate.RegisterValidation("notallzero", common.NotAllZero)
	}

	db := database.DatabaseInit(common.GetDbURL())
	log.Println("Database connection successful")
	dbHandler := handler.NewHandler(db)
	api := router.NewRouter(dbHandler)
	log.Println("Starting the API")
	api.StartAPI()
}
