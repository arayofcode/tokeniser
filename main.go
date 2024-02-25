package main

import (
	"log"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
)

func main() {
	db := database.DatabaseInit(common.GetDbURL())
	log.Println("Database connection successful")
	dbHandler := handler.NewHandler(db)
	api := router.NewRouter(dbHandler)
	log.Println("Starting the API")
	api.StartAPI()
}

/*
Possible Idempotency:
- Request with same payload sent again
- Same key-value pair sent again. Should we generate new tokens in this case? Think about same names of two different people

User --> API --> Handler --> Database

Next steps:
- Data Validation
- Logging
- Encryption
*/
