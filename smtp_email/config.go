package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type SDBConfig struct {
	Host     string
	Password string
	Name     string
	Port     string
	SSLMode  string
	TimeZone string
}

type SSMTPConfig struct {
	Host               string
	Port               string
	Sender             string
	Username           string
	Password           string
	InsecureSkipVerify bool
	StartTLS           bool
	IsAuthRequired     bool
}

var (
	DB_CONFG    SDBConfig
	SMTP_CONFIG SSMTPConfig
)

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// Read database configuration
	dbConfig := viper.GetStringMapString("database")
	DB_CONFG = SDBConfig{
		Host:     dbConfig["host"],
		Password: dbConfig["password"],
		Name:     dbConfig["name"],
		Port:     dbConfig["port"],
		SSLMode:  dbConfig["sslmode"],
		TimeZone: dbConfig["timezone"],
	}
	smtpConfig := viper.GetStringMapString("smtp")

	insecureSkipVerifyStr := smtpConfig["insecure_skip_verify"]
	insecureSkipVerify, err := strconv.ParseBool(insecureSkipVerifyStr)
	if err != nil {
		fmt.Println("insecureSkipVerify: Error parsing boolean:", err)
		return
	}
	startTlsStr := smtpConfig["start_tls"]
	startTls, err := strconv.ParseBool(startTlsStr)
	if err != nil {
		fmt.Println("startTls: Error parsing boolean:", err)
		return
	}

	isAuthRequiredStr := smtpConfig["is_auth_required"]
	isAuthRequired, err := strconv.ParseBool(isAuthRequiredStr)
	if err != nil {
		fmt.Println("isAuthRequired: Error parsing boolean:", err)
		return
	}

	SMTP_CONFIG = SSMTPConfig{
		Host:               smtpConfig["host"],
		Port:               smtpConfig["port"],
		Sender:             smtpConfig["sender"],
		Username:           smtpConfig["username"],
		Password:           smtpConfig["password"],
		InsecureSkipVerify: insecureSkipVerify,
		StartTLS:           startTls,
		IsAuthRequired:     isAuthRequired,
	}
}

var EMAIL_LIST = []string{
	// "test@billowdev.com",
	"a1@billowdev.com",
	// "a2@billowdev.com",
	// "a9@gmail.com",
}
