package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/arayofcode/tokeniser/cipher"
	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/router"
)

func init() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("expirydate", common.ExpiryDateMMYY)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to register 'expirydate' validation")
		}

		err = validate.RegisterValidation("notallzero", common.NotAllZero)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to register 'notallzero' validation")
		}
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbURL := common.GetDbURL()
	if dbURL == "" {
		log.Fatal().Msg("Database URL is not configured")
	}

	db := database.DatabaseInit(dbURL)
	log.Info().Msg("Database connection successful")

	passphrase := common.GetPassphrase()
	if passphrase == "" {
		log.Fatal().Msg("Encryption passphrase is not configured")
	}

	cipherModule := cipher.Init(passphrase)

	dbHandler := handler.NewHandler(db)
	apiRouter := router.NewRouter(dbHandler, cipherModule)

	log.Info().Msg("Starting the API")
	apiRouter.StartAPI()
}
