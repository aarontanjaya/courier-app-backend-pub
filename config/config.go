package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type authConfig struct {
	Secret   string
	TokenDur string
}

type referralConfig struct {
	ClaimerReward   float64
	ClaimerTreshold float64
	OwnerReward     float64
	OwnerTreshold   float64
}

type topupConfig struct {
	TopupMin float64
	TopupMax float64
}

type AppConfig struct {
	DBConfig       dbConfig
	AuthConfig     authConfig
	ReferralConfig referralConfig
	TopupConfig    topupConfig
	ENV            string
}

func getENV(key, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}

var Config AppConfig

func InitENV() {
	godotenv.Load()
	claimerReward, err := strconv.ParseFloat(os.Getenv("REFERRAL_CLAIMER_REWARD"), 64)
	if err != nil {
		panic("REFERRAL_CLAIMER_REWARD not set in env")
	}

	claimerTreshold, err := strconv.ParseFloat(os.Getenv("REFERRAL_CLAIMER_TRESHOLD"), 64)
	if err != nil {
		panic("REFERRAL_CLAIMER_TRESHOLD not set in env")
	}

	ownerReward, err := strconv.ParseFloat(os.Getenv("REFERRAL_OWNER_REWARD"), 64)
	if err != nil {
		panic("REFERRAL_OWNER_REWARD not set in env")
	}

	ownerTreshold, err := strconv.ParseFloat(os.Getenv("REFERRAL_OWNER_TRESHOLD"), 64)
	if err != nil {
		panic("REFERRAL_OWNER_TRESHOLD not set in env")
	}

	topupMinTreshold, err := strconv.ParseFloat(os.Getenv("TOPUP_MIN_TRESHOLD"), 64)
	if err != nil {
		panic("TOPUP_MIN_TRESHOLD not set in env")
	}
	topupMaxTreshold, err := strconv.ParseFloat(os.Getenv("TOPUP_MAX_TRESHOLD"), 64)
	if err != nil {
		panic("TOPUP_MAX_TRESHOLD not set in env")
	}
	Config = AppConfig{
		DBConfig: dbConfig{
			Host:     getENV("DB_HOST", ""),
			User:     getENV("DB_USER", ""),
			Password: getENV("DB_PASSWORD", ""),
			DBName:   getENV("DB_NAME", ""),
			Port:     getENV("DB_PORT", ""),
		},
		AuthConfig: authConfig{
			Secret:   getENV("SECRET", ""),
			TokenDur: getENV("TOKEN_DUR", "900"),
		},
		ReferralConfig: referralConfig{
			ClaimerReward:   claimerReward,
			ClaimerTreshold: claimerTreshold,
			OwnerReward:     ownerReward,
			OwnerTreshold:   ownerTreshold,
		},
		TopupConfig: topupConfig{
			TopupMin: topupMinTreshold,
			TopupMax: topupMaxTreshold,
		},
		ENV: getENV("ENV", "production"),
	}
}
