package models

type Config struct {
	DB         string `json:"db_creds"`
	Passphrase string `json:"passphrase"`
}
