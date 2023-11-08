package helper

import (
	"encoding/hex"
	"math/rand"
	"time"
	"unicode"
)

func GenerateAPIKey() string {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return err.Error()
	}

	apiKey := hex.EncodeToString(seed)

	return apiKey
}

func GeneratePassword() string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-=[]{}|<>/?~"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func ValidatePasswordStrength(password string) bool {
	var (
		upp, low, num, sym bool
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			low = true
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsNumber(char):
			num = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym {
		return false
	}

	return true
}

func ParseDate(layout, date string) error {
	_, err := time.Parse(layout, date)
	if err != nil {
		return err
	}

	return nil
}

func FormatStartTimeForSQL(date string) string {
	return date + " 00:00:00"
}

func FormatEndTimeForSQL(date string) string {
	return date + " 24:00:00"
}

func FormatWIB(currentTime time.Time) string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime = currentTime.In(loc)

	return currentTime.Format("2006-01-02 15:04:05 MST")
}
