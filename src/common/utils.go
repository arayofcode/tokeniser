package common

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/arayofcode/tokeniser/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
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
	cfg.DB = os.Getenv("DB")
	cfg.Passphrase = os.Getenv("PASSPHRASE")
	return
}

// func loadConfig() models.Config {
// 	configPath := os.Getenv("CONFIG_PATH")
// 	if configPath == "" {
// 		log.Fatal().Msg("CONFIG_PATH environment variable is not set")
// 	}
//
// 	jsonFile, err := os.Open(configPath)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error opening config file")
// 	}
// 	defer jsonFile.Close()
//
// 	byteValue, err := io.ReadAll(jsonFile)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error reading config file")
// 	}
//
// 	var cfg models.Config
// 	if err := json.Unmarshal(byteValue, &cfg); err != nil {
// 		log.Fatal().Err(err).Msg("Error unmarshalling config")
// 	}
//
// 	return cfg
// }

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
