package common

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/arayofcode/tokeniser/src/models"
)

var config *models.Config
var once sync.Once

func PrettyPrint(data interface{}) string {
	var p []byte
	p, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling data for pretty print")
		return ""
	}
	return string(p)
}

func loadConfig() (cfg models.Config) {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("%+v\n", err)
	}
	cfg.DB = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=10", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	return cfg
}

func getConfig() models.Config {
	once.Do(func() {
		config = new(models.Config)
		*config = loadConfig()
	})
	return *config
}

func GetDbURL() string {
	return getConfig().DB
}

func GetPassphrase() string {
	return getConfig().Passphrase
}

func GetPort() string {
	return getConfig().AppPort
}

var ExpiryDateMMYY validator.Func = func(fl validator.FieldLevel) bool {
	expiry := fl.Field().String()
	if len(expiry) != 4 {
		return false
	}

	month, err := strconv.Atoi(expiry[:2])
	if err != nil || month < 1 || month > 12 {
		return false
	}

	year, err := strconv.Atoi(expiry[2:])
	if err != nil {
		return false
	}

	currentYear := time.Now().Year() % 100
	currentMonth := int(time.Now().Month())

	if year < currentYear || (year == currentYear && month < currentMonth) {
		return false
	}

	return true
}

var NotAllZero validator.Func = func(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	isValidLength := len(cardNumber) >= 13 && len(cardNumber) <= 19
	isAllZeroes := regexp.MustCompile(`^0+$`).MatchString(cardNumber)
	return isValidLength && !isAllZeroes
}

func AssertByteSliceEqual(t *testing.T, expected, actual []byte) {
	if (expected == nil && len(actual) == 0) || (actual == nil && len(expected) == 0) {
		assert.True(t, true)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func MaskLeft(s string, lastVisibleN int) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-lastVisibleN; i++ {
		rs[i] = 'X'
	}
	return string(rs)
}
