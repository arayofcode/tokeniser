package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
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

func ValidateExpiryDateMMYYYY(fl validator.FieldLevel) bool {
	expiryDate := fl.Field().String()
	_, err := time.Parse("01/2006", expiryDate)
	return err == nil
}

func ValidateNotAllZero(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	// Check for length between 16 and 19 and ensure it's not all zeros
	isValidLength := len(cardNumber) >= 16 && len(cardNumber) <= 19
	isAllZeroes := regexp.MustCompile(`^0+$`).MatchString(cardNumber)
	return isValidLength && !isAllZeroes
}
