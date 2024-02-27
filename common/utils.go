package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func PrettyPrint(data interface{}) string {
	var p []byte
	p, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ""
	}
	return string(p)
}

func GetDbURL() (url string) {
	jsonFile, err := os.Open("/Users/aryansharma/code/learning/tokeniser/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var config struct {
		Db string `json:"db_creds"`
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config file: %v\n", err)
		return
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing json: %v\n", err)
		return
	}

	return config.Db
}

func ExpiryDateMMYY(fl validator.FieldLevel) bool {
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

	// Expiry date is a previous date
	if year < currentYear || (year == currentYear && month < currentMonth) {
		return false
	}

	return true
}

func NotAllZero(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	isValidLength := len(cardNumber) >= 13 && len(cardNumber) <= 19
	isAllZeroes := regexp.MustCompile(`^0+$`).MatchString(cardNumber)
	return isValidLength && !isAllZeroes
}
