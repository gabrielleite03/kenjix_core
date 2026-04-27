package config

import "os"

type Config struct {
	ClientID     string
	ClientSecret string
	SellerID     string
	DBUrl        string
}

func Load() *Config {
	return &Config{
		ClientID:     os.Getenv("ML_CLIENT_ID"),
		ClientSecret: os.Getenv("ML_CLIENT_SECRET"),
		SellerID:     os.Getenv("ML_SELLER_ID"),
		DBUrl:        os.Getenv("DATABASE_URL"),
	}
}
