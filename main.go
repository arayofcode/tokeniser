package main

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/arayofcode/tokeniser/cipher"
	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
)

func init() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("expirydate", common.ExpiryDateMMYY)
		validate.RegisterValidation("notallzero", common.NotAllZero)
	}
}

func main() {
	db := database.DatabaseInit(common.GetDbURL())
	log.Println("Database connection successful")
	dbHandler := handler.NewHandler(db)
	cipher := cipher.Init(common.GetPassphrase())
	api := router.NewRouter(dbHandler, cipher)
	log.Println("Starting the API")
	api.StartAPI()
}
