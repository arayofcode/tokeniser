package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func PrettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Printf("%s \n\n", p)
}

func GetDbURL() (url string) {
	jsonFile, err := os.Open("config.json")
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
