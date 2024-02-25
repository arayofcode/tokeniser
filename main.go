package main

import (
	"log"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
	"github.com/go-playground/validator/v10"
)

func main() {
	// registering validator
	validate := validator.New()
	validate.RegisterValidation("expiry_date", common.ValidateExpiryDateMMYYYY)
	validate.RegisterValidation("notallzero", common.ValidateNotAllZero)

	db := database.DatabaseInit(common.GetDbURL())
	log.Println("Database connection successful")
	dbHandler := handler.NewHandler(db)
	api := router.NewRouter(dbHandler)
	log.Println("Starting the API")
	api.StartAPI()
}
