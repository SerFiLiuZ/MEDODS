package apiserver

import (
	"os"

	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DBconnecturi string
	Port         string
	JwtKey       []byte
}

const (
	envFilePath string = "../../config/.env"
)

func LoadEnv(logger *utils.Logger) error {
	logger.Debugf("Loading .env file from path: %s", envFilePath)

	err := godotenv.Load(envFilePath)
	if err != nil {
		logger.Fatal("Error loading .env file: %v", err)
		return err
	}

	return nil
}

func GetConfig() *Config {
	return &Config{
		DBconnecturi: os.Getenv("DB_CONNECT_URI"),
		Port:         os.Getenv("PORT"),
		JwtKey:       []byte(os.Getenv("JWTKEY")),
	}
}
