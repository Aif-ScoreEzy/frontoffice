package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	Env                            string
	FrontendBaseUrl                string
	Port                           string
	DbUser                         string
	DbPassword                     string
	DbName                         string
	DbPort                         string
	DbHost                         string
	MailjetEmail                   string
	MailtjetUsername               string
	MailjetPublicKey               string
	MailjetSecretKey               string
	JwtSecretKey                   string
	JwtExpiresMinutes              string
	JwtVerificationExpiresMinutes  string
	JwtActivationExpiresMinutes    string
	JwtResetPasswordExpiresMinutes string
	AifcoreHost                    string
	GenretailV3                    string
	GetLogByDate                   string
	GetLogByRangeDate              string
	GetLogByMonth                  string
	GetLogByName                   string
}

func GetEnvironment(key string) string {
	return os.Getenv(key)
}

func LoadEnvironment() *Environment {
	env := os.Getenv("APP_ENV")
	if env == "" {
		fmt.Println("No App env")
		env = "local"
	}

	if env == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error loading .env file")
		}
	}

	return &Environment{
		Env:                            GetEnvironment("APP_ENV"),
		FrontendBaseUrl:                GetEnvironment("FRONTEND_BASE_URL"),
		Port:                           GetEnvironment("APP_PORT"),
		DbUser:                         GetEnvironment("DB_USER"),
		DbPassword:                     GetEnvironment("DB_PASSWORD"),
		DbName:                         GetEnvironment("DB_NAME"),
		DbPort:                         GetEnvironment("DB_PORT"),
		DbHost:                         GetEnvironment("DB_HOST"),
		MailjetEmail:                   GetEnvironment("MAILJET_EMAIL"),
		MailtjetUsername:               GetEnvironment("MAILJET_USERNAME"),
		MailjetPublicKey:               GetEnvironment("MAILJET_PUBLIC_KEY"),
		MailjetSecretKey:               GetEnvironment("MAILJET_SECRET_KEY"),
		JwtSecretKey:                   GetEnvironment("JWT_SECRET_KEY"),
		JwtExpiresMinutes:              GetEnvironment("JWT_EXPIRES_MINUTES"),
		JwtVerificationExpiresMinutes:  GetEnvironment("JWT_VERIFICATION_EXPIRES_MINUTES"),
		JwtActivationExpiresMinutes:    GetEnvironment("JWT_ACTIVATION_EXPIRES_MINUTES"),
		JwtResetPasswordExpiresMinutes: GetEnvironment("JWT_RESET_PASSWORD_EXPIRES_MINUTES"),
		AifcoreHost:                    GetEnvironment("AIFCORE_HOST"),
		GenretailV3:                    GetEnvironment("GEN_RETAIL_V3"),
		GetLogByDate:                   GetEnvironment("GET_LOGS_BY_DATE"),
		GetLogByRangeDate:              GetEnvironment("GET_LOGS_BY_RANGE_DATE"),
		GetLogByMonth:                  GetEnvironment("GET_LOGS_BY_MONTH"),
		GetLogByName:                   GetEnvironment("GET_LOGS_BY_NAME"),
	}
}
