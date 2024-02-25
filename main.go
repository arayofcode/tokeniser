package main

import (
	"context"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
)

func main() {
	ctx := context.Background()
	db := database.DatabaseInit(common.GetDbURL())
	dbHandler := handler.NewHandler(db)
	api := router.NewRouter(ctx, dbHandler)
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
